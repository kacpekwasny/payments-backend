package funcs

import (
	"encoding/json"
	"fmt"
	"net/http"

	cmt "github.com/kacpekwasny/commontools"
)

/*
//  ok
//  internal_error
//  login_not_found
//  pass_missmatch
//  login_in_use
//  login_requirements
//  pass_requirements
//  paswd_pwned
//  invalid_chars
//  unauth
//  room_link_not_found
//  no_rows  */
func Respond(w http.ResponseWriter, msg_title string, key_val ...interface{}) {

	ln := len(key_val)
	if ln%2 != 0 {
		panic("len of key_val has to be event")
	}

	err_code := GetErrCode(msg_title)
	w.Header().Add("Content-Type", "application/json")
	resp := map[string]interface{}{
		"err_code": err_code,
	}
	for i := 0; i < ln; i += 2 {
		resp[key_val[i].(string)] = key_val[i+1]
	}
	cmt.Pc(json.NewEncoder(w).Encode(resp))
}

var ErrCodes = map[string]int{
	"ok":                  0,
	"internal_error":      1,
	"login_not_found":     2,
	"pass_missmatch":      3,
	"login_in_use":        5,
	"login_requirements":  6,
	"pass_requirements":   7,
	"paswd_pwned":         8,
	"invalid_chars":       9,
	"unauth":              10,
	"room_link_not_found": 11,
	"no_rows":             13,
}

func GetErrCode(message_title string) int {
	code, ok := ErrCodes[message_title]
	if !ok {
		fmt.Printf("ErrCodes[ %v ] is missing \n", message_title)
		return 1
	}
	return code
}

// RIE(w) == Respond(w, "internal_error")
func RIE(w http.ResponseWriter) {
	Respond(w, "internal_error")
}
