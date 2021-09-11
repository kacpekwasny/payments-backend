package router

import "net/http"

func updateAcceptPayment(w http.ResponseWriter, r *http.Request)   {}
func updateUnacceptPayment(w http.ResponseWriter, r *http.Request) {}
func updateChangeUserRole(w http.ResponseWriter, r *http.Request)  {}

var (
	UpdateAcceptPayment   = http.HandlerFunc(updateAcceptPayment)
	UpdateUnacceptPayment = http.HandlerFunc(updateUnacceptPayment)
	UpdateChangeUserRole  = http.HandlerFunc(updateChangeUserRole)
)
