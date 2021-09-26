package api_router

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kacpekwasny/payments-backend/internal/funcs"
)

func addUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	uname := params["username"]
	_, err := Ctm.Exec(`INSERT INTO payments.users (username) VALUES (?)`,
		[]interface{}{uname})
	if err != nil {
		fmt.Println(err)
		funcs.RIE(w)
		return
	}
	funcs.ROK(w)
}

var (
	AddUser = http.HandlerFunc(addUser)
)
