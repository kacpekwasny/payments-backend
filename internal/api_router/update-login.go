package api_router

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	f "github.com/kacpekwasny/payments-backend/internal/funcs"
)

func updateLogin(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	oldLogin := params["old_login"]
	newLogin := params["new_login"]
	fmt.Println("updateLogin()", oldLogin, " -> ", newLogin)
	res, err := Ctm.Exec(`
		UPDATE payments.users
		SET username=?
		WHERE username=?`,
		[]interface{}{newLogin, oldLogin})
	i, errRows := res.RowsAffected()
	fmt.Println("updateLogin() res.RowsAffected()", i, errRows)
	if err != nil {
		fmt.Println("updateLogin: ", err)
		f.RIE(w)
		return
	}
	f.ROK(w)
}

var (
	UpdateLogin = http.HandlerFunc(updateLogin)
)
