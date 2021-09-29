package wrappers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kacpekwasny/payments-backend/internal/funcs"
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
		resp, err := client.Get(Config.AuthApiBaseUrl + fmt.Sprintf("/isAuthenticated/%s/%s", uname, token))
		if err != nil {
			funcs.RIE(w)
			fmt.Println("client.Get error:", err)
			return
		}

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			funcs.RIE(w)
			fmt.Println("io.ReadAll error:", err)
			return
		}
		m := map[string]int{}
		err = json.Unmarshal(b, &m)
		if err != nil {
			funcs.RIE(w)
			fmt.Println("json.Unmarshal error:", err)
			return
		}

		fmt.Println(m)
		if m["err_code"] != 0 {
			// Unauth
			funcs.Respond(w, "unauth")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func AuthorisedForRoom(next http.Handler) http.Handler {
	// This function
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get username
		params := mux.Vars(r)
		uname := params["username"]
		roomLink := params["room_link"]
		fmt.Println(uname, roomLink)

		row, err := Ctm.QueryRow("CALL AuthorisedForRoom(?, ?)",
			[]interface{}{roomLink, uname})

		if err != nil {
			funcs.RIE(w)
			return
		}
		found := 0
		err = row.Scan(&found)
		if err != nil {
			fmt.Println("AuthorisedForRooom, Row.Scan, err:", err)
			funcs.RIE(w)
			return
		}
		if found != 1 {
			funcs.Respond(w, "unauth")
			return
		}

		// User is authorised for this room
		fmt.Printf("User '%s' authorised for room of link '%s' \n", uname, roomLink)
		next.ServeHTTP(w, r)
	})
}

func UserIsAdminForRoom(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		uname := params["username"]
		roomLink := params["room_link"]

		row, err := Ctm.QueryRow(`SELECT is_admin
			FROM payments.rooms_have_users
			WHERE rooms_link=? AND users_id=(SELECT id FROM users WHERE username=?)`,
			[]interface{}{roomLink, uname})

		if err != nil {
			funcs.RIE(w)
			fmt.Println(err)
			if err != sql.ErrNoRows {
				panic(err)
			}
			return
		}
		isAdmin := false
		row.Scan(&isAdmin)
		if !isAdmin {
			fmt.Println("non Admin tried accessing admin func")
			funcs.Respond(w, "not_admin")
			return
		}

		next.ServeHTTP(w, r)
	})
}

// User is authenticated and authorised
func AA(next http.Handler) http.Handler {
	if Config.AuthON {
		return UserIsAuthenticated(
			AuthorisedForRoom(
				next))

	}
	fmt.Println("AA wrapper is disabled for development")
	return next
}

func AAU(next http.Handler) http.Handler {
	return UserIsAuthenticated(
		AuthorisedForRoom(
			UserIsAdminForRoom(next)))
}
