package router

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	f "github.com/kacpekwasny/payments-backend/internal/funcs"
	jstr "github.com/kacpekwasny/payments-backend/internal/json_structs"
	scm "github.com/kacpekwasny/payments-backend/pkg/sql_conn_manager"
)

func getRoomData(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ch := make(chan scm.QueryRet)
	// get name and desc by link
	Ctm.StmtChan <- scm.InStmt{
		SqlText:    "SELECT name, desc FROM payments_room.rooms WHERE link=?",
		StmtArgs:   []interface{}{params["room_link"]},
		ActionType: 2,
		ResChan:    ch,
	}

	qret := <-ch
	if qret.Err != nil {
		fmt.Println("getRoomData, qret.Err:", qret.Err)
		fmt.Fprint(w, "Internal error")
		return
	}
	var title, desc string
	err := qret.Row.Scan(&title, &desc)
	if err != nil {
		f.RIE(w)
		panic(err)
	}
	f.Respond(w, "ok",
		"title", title,
		"desc", desc)

}
func getRoomPayments(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ch := make(chan scm.QueryRet)

	Ctm.StmtChan <- scm.InStmt{
		SqlText: `SELECT pft.id as 'paymemnt_id', pft.title, pft.value, created, from_accepted, to_accepted,
		(SELECT username FROM users WHERE pft.from_user_id=users.id) as 'from',
		(SELECT username FROM users WHERE pft.to_user_id=users.id) as 'to'
		FROM payments_from_to as pft
		WHERE pft.room_id = (SELECT id FROM rooms WHERE link=?)
		ORDER BY pft.created ASC`,
		StmtArgs:   []interface{}{params["room_link"]},
		ActionType: 1,
		ResChan:    ch,
	}

	qret := <-ch
	if qret.Err != nil {
		if qret.Err == sql.ErrNoRows {
			f.Respond(w, "no_rows")
			return
		}
		f.RIE(w)
		panic(qret.Err)
	}

	pftLs := []jstr.PaymentFromTo{}
	for qret.Rows.Next() {
		pft := jstr.PaymentFromTo{}
		err := qret.Rows.Scan(&pft.ID, &pft.Title, &pft.Value, &pft.Created, &pft.FromAccepted,
			&pft.ToAccepted, &pft.FromUsername, &pft.ToUsername)
		if err != nil {
			f.RIE(w)
			panic(err)
		}
		pftLs = append(pftLs, pft)
	}

	f.Respond(w, "ok",
		"payments_from_to", pftLs)
}

var GetRoomData = http.HandlerFunc(getRoomData)
var GetRoomPayments = http.HandlerFunc(getRoomPayments)
