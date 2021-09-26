package router

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kacpekwasny/payments-backend/internal/funcs"
)

// accept or unaccept the payment
func updateSetAcceptPayment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	uname := params["username"]
	paymentID := params["payment_id"]

	accept, ok := map[string]int{
		"unaccept": 0,
		"accept":   1,
	}[params["accept"]]

	if !ok {
		funcs.Respond(w, "incorrect_input",
			"accept_can_be_only", "accept or unaccept")
		fmt.Println("updateSetAcceptPayment, param accept, wrong value.")
		return
	}

	// we dont know if the user is accepting user payment or event payment
	// so first we try accepting user payment,
	_, err := Ctm.Exec(`CALL SetAcceptUserPayment(?, ?, ?)`,
		[]interface{}{uname, accept, paymentID})

	if err != nil {
		funcs.RIE(w)
		panic(err)
		return
	}

	// and then event payment
	rows, err := Ctm.Query(`CALL SetAcceptEventPayment(?, ?, ?)`,
		[]interface{}{uname, accept, paymentID})

	if err != nil {
		funcs.RIE(w)
		panic(err)
		return
	}

	// if the procedure returned rows, handle them
	handleRowsEventPayment(paymentID, rows)

	funcs.ROK(w)
}

func updateChangeUserRole(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userForRoleChange := params["user_for_change"]
	userExecutingChange := params["username"]
	if userForRoleChange == userExecutingChange {
		funcs.Respond(w, "incorrect_input",
			"msg", "cannot change your own role")
		return
	}
	roomLink := params["room_link"]
	newRole, ok := map[string]int{
		"standard": 0,
		"admin":    1,
	}[params["new_role"]]

	if !ok {
		funcs.Respond(w, "incorrect_input",
			"new_role_can_be_only", "standard or admin")
		fmt.Println("updateChangeUserRole, param new_role, wrong value.")
		return
	}

	result, err := Ctm.Exec(`UPDATE payments.rooms_have_users
		SET is_admin=?
		WHERE users_id=(SELECT id FROM payments.users WHERE username=?)
		AND rooms_link=?`,
		[]interface{}{newRole, userForRoleChange, roomLink})

	if err != nil {
		funcs.RIE(w)
		panic(err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		funcs.ROK(w)
		panic(err)
	}

	if affected == 0 {
		fmt.Println("No rows affected", userForRoleChange, params["new_role"])
		funcs.Respond(w, "ok",
			"msg", "no rows were affected")
		return
	}
	funcs.ROK(w)
}

func updateRoomNameAndDescription(w http.ResponseWriter, r *http.Request) {
	// require admin
	params := mux.Vars(r)
	roomLink := params["room_link"]
	roomName, roomDesc, err := funcs.PrepareNameDescFromRequest(r)
	if err != nil {
		fmt.Println(err)
		funcs.RIE(w)
		return
	}
	changedName := false
	changedDesc := false
	if 0 < len(roomName) && len(roomName) < 250 {
		changedName = true
		_, err := Ctm.Exec(`
			UPDATE payments.rooms
			SET name=?
			WHERE link=?`,
			[]interface{}{roomName, roomLink})
		if err != nil {
			changedName = false
			fmt.Println(err)
		}
	}

	if 0 < len(roomDesc) && len(roomDesc) < 10000 {
		changedDesc = true
		_, err := Ctm.Exec(`
			UPDATE payments.rooms
			SET `+"`desc`"+`=?
			WHERE link=?`,
			[]interface{}{roomDesc, roomLink})
		if err != nil {
			changedDesc = false
			fmt.Println(err)
		}
	}
	funcs.Respond(w, "ok",
		"changed_name", changedName,
		"changed_desc", changedDesc)
}

var (
	UpdateSetAcceptPayment       = http.HandlerFunc(updateSetAcceptPayment)
	UpdateChangeUserRole         = http.HandlerFunc(updateChangeUserRole)
	UpdateRoomNameAndDescription = http.HandlerFunc(updateRoomNameAndDescription)
)

func handleRowsEventPayment(paymentID string, rows *sql.Rows) {
	// procedure 'SetAcceptEventPayment' in the database will return rows of {username, input}
	// where input is the ammount user paid for the event
	// BUT only if all users accepted this payment.
	sum := 0.0
	usersNum := 0
	unameAmm := map[string]float64{}

	for rows.Next() {
		uname := ""
		var ammount float64
		err := rows.Scan(&uname, &ammount)
		if err != nil {
			panic(err)
		}
		unameAmm[uname] = ammount
		sum += ammount
		usersNum++
	}
	if usersNum == 0 {
		return
	}

	// get room link
	row, err := Ctm.QueryRow(`
		SELECT rooms_link
		FROM payments.event_payments
		WHERE public_id=?`,
		[]interface{}{paymentID})
	if err != nil {
		fmt.Println(err)
		return
	}
	roomLink := ""
	err = row.Scan(&roomLink)
	if err != nil {
		fmt.Println(err)
		return
	}
	// The calculated delta should be applied to saldo in database
	// SQL: for every user: saldo = saldo + delta
	avg := sum / float64(usersNum)
	for uname, amm := range unameAmm {
		delta := amm - avg

		go func(delta float64, roomLink, uname string) {
			result, err := Ctm.Exec(`
				UPDATE payments.rooms_have_users
				SET saldo=saldo+?
				WHERE 
					rooms_link=?
				AND
					users_id=(
						SELECT id
						FROM payments.users
						WHERE username=?
					);`,
				[]interface{}{delta, roomLink, uname})
			if err != nil {
				fmt.Printf("UPDATE saldo %#v, %#v, %#v \n  %#v, %#v", delta, roomLink, uname, result, err)
			}
		}(delta, roomLink, uname)
	}
}
