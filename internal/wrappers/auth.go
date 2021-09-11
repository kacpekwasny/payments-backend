package wrappers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kacpekwasny/payments-backend/internal/funcs"
	scm "github.com/kacpekwasny/payments-backend/pkg/sql_conn_manager"
)

func UserIsAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		uname := params["username"]
		token := params["token"]

		// User input! Let's make sure it is safe
		if !(funcs.NoForbidenChars(uname) && funcs.NoForbidenChars(token)) {
			// forbiden chars found! no response for U!
			fmt.Println("FORBIDEN CHARS: UserIsAuthenticated: username:", uname, "token:", token)
			return
		}

		// input is safe
		resp, err := client.Get(config.AuthApiBaseUrl + fmt.Sprintf("/%s/%s", uname, token))
		if err != nil {
			fmt.Println(err)
			fmt.Fprint(w, "Internal fail")
		}

		b, err := io.ReadAll(resp.Body)
		m := map[string]int{}
		json.Unmarshal(b, &m)

		if m["err_code"] == 0 {
			next.ServeHTTP(w, r)
		}

		// Unauth
		fmt.Fprint(w, "unauth")
	})
}

func AuthorisedForRoom(next http.Handler) http.Handler {
	// This function
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get username
		params := mux.Vars(r)
		uname := params["username"]
		roomId := params["room_id"]

		resChan := make(chan scm.QueryRet)
		defer close(resChan)
		Ctm.StmtChan <- scm.InStmt{
			SqlText:    "CALL AuthorisedForRoom(?, ?)",
			StmtArgs:   []interface{}{uname, roomId},
			ActionType: 2,
			ResChan:    resChan,
		}

		queryRet := <-resChan
		if queryRet.Err != nil {
			fmt.Fprint(w, "Internal error")
			return
		}
		found := 0
		err := queryRet.Row.Scan(&found)
		if err != nil {
			fmt.Println("AuthorisedForRooom, Row.Scan, err:", err)
			fmt.Fprint(w, "Internal error")
			return
		}
		if found != 1 {
			fmt.Fprint(w, "Unauthorised for room")
			return
		}

		// User is authorised for this room
		next.ServeHTTP(w, r)
	})
}

// User is authenticated and authorised
func AA(next http.Handler) http.Handler {
	return UserIsAuthenticated(
		AuthorisedForRoom(
			next))
}
