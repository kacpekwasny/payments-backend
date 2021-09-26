package router

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kacpekwasny/payments-backend/internal/funcs"
	f "github.com/kacpekwasny/payments-backend/internal/funcs"
	jstr "github.com/kacpekwasny/payments-backend/internal/json_structs"
)

func postNewPayment(w http.ResponseWriter, r *http.Request) {
	pay, err := funcs.PreparePaymentFromRequest(r)
	if err != nil {
		funcs.Respond(w, "incorrect_input")
		fmt.Println("routes-post.go, postNewPayment", err)
		return
	}
	params := mux.Vars(r)
	uname := params["username"]
	roomLink := params["room_link"]

	ok := false
	if pay.Type == jstr.USER_PAYMENT {
		if !jstr.UserPaymentIsValid(pay) {
			fmt.Println("!jstr.UserPaymentIsValid(pay)")
			return
		}
		ok = insertUserPayment(roomLink, uname, pay)
	}

	if pay.Type == jstr.EVENT_PAYMENT {
		if !jstr.EventPaymentIsValid(pay) {
			fmt.Println("!jstr.EventPaymentIsValid(pay)")
			return
		}
		ok = insertEventPayment(roomLink, uname, pay)
	}

	if ok {
		funcs.ROK(w)
		return
	}

	fmt.Println("insertPayment not OK")
	funcs.RIE(w)
}

func postAddUserRoom(w http.ResponseWriter, r *http.Request) {
	// requires admin, appropriate wrapper should be applied
	params := mux.Vars(r)
	roomLink := params["room_link"]
	newUser := params["new_user"]

	_, err := Ctm.Exec(`
	INSERT INTO payments.rooms_have_users
		(rooms_id, rooms_link, users_id)
	VALUES (
		(SELECT id FROM payments.rooms WHERE link=?),
		?,
		(SELECT id FROM payments.users WHERE username=?)
	)`,
		[]interface{}{roomLink, roomLink, newUser})
	if err != nil {
		funcs.RIE(w)
		fmt.Println(err)
		return
	}
	funcs.ROK(w)
}

func postRemoveUserRoom(w http.ResponseWriter, r *http.Request) {
	// requires admin, appropriate wrapper should be applied
	params := mux.Vars(r)
	roomLink := params["room_link"]
	uToRemove := params["old_user"]

	_, err := Ctm.Exec(`
	DELETE FROM payments.rooms_have_users
	WHERE
		rooms_link=?
	  AND
		users_id=(
			SELECT id
			FROM payments.users
			WHERE username=?
		)`,
		[]interface{}{roomLink, uToRemove})
	if err != nil {
		funcs.RIE(w)
		fmt.Println(err)
		return
	}
	funcs.ROK(w)
}

func postNewRoom(w http.ResponseWriter, r *http.Request) {
	// requires authentication NOT authorization
	params := mux.Vars(r)
	uname := params["username"]

	row, err := Ctm.QueryRow("CALL InsertRoom(?, ?, ?)",
		[]interface{}{"Title of new room", "Description of new room", uname})
	if err != nil {
		fmt.Println(err)
		f.RIE(w)
		panic(err)
		return
	}

	roomLink := ""
	err = row.Scan(&roomLink)
	if err != nil {
		fmt.Println(err)
		f.RIE(w)
		panic(err)
		return
	}

	f.Respond(w, "ok",
		"link", roomLink)
}

func leaveRoom(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	uname := params["username"]
	roomLink := params["room_link"]

	_, err := Ctm.Exec(`CALL LeaveRoom(?, ?)`,
		[]interface{}{uname, roomLink})
	if err != nil {
		fmt.Println("leaveRoom", err)
		f.RIE(w)
		return
	}
	f.ROK(w)
}

var (
	PostNewPayment     = http.HandlerFunc(postNewPayment)
	PostAddUserRoom    = http.HandlerFunc(postAddUserRoom)
	PostRemoveUserRoom = http.HandlerFunc(postRemoveUserRoom)
	PostNewRoom        = http.HandlerFunc(postNewRoom)
	LeaveRoom          = http.HandlerFunc(leaveRoom)
)
