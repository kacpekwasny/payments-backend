package router

import "net/http"

func postNewPayment(w http.ResponseWriter, r *http.Request)     {}
func postAddUserRoom(w http.ResponseWriter, r *http.Request)    {}
func postRemoveUserRoom(w http.ResponseWriter, r *http.Request) {}

var (
	PostNewPayment     = http.HandlerFunc(postNewPayment)
	PostAddUserRoom    = http.HandlerFunc(postAddUserRoom)
	PostRemoveUserRoom = http.HandlerFunc(postRemoveUserRoom)
)
