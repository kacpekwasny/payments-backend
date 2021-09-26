package router

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	f "github.com/kacpekwasny/payments-backend/internal/funcs"
)

func getRoomData(w http.ResponseWriter, r *http.Request) {
	// previous functions allready used params which we will use
	params := mux.Vars(r)
	chInfo := make(chan sqlGetInfo, 1)
	chUserPay := make(chan sqlGetPay, 1)
	chEvent := make(chan sqlGetPay, 1)
	chUsers := make(chan usersInfo, 1)
	chOrder := make(chan []string, 1)
	defer close(chInfo)
	defer close(chUserPay)
	defer close(chEvent)
	defer close(chUsers)
	defer close(chOrder)

	var link = params["room_link"]
	go sqlGetRoomInfo(chInfo, link, params["username"])
	go sqlGetUserPayments(chUserPay, link)
	go sqlGetEventPayments(chEvent, link)
	go sqlGetRoomUsers(chUsers, params["room_link"])
	go sqlGetPaymentsOrder(chOrder, link)

	info := <-chInfo
	userPay := <-chUserPay
	eventPay := <-chEvent
	users := <-chUsers
	order := <-chOrder

	if info.err || userPay.err || eventPay.err || users.err || order == nil {
		f.RIE(w)
		fmt.Println("ERROR: info.err || userPay.err || eventPay.err || users.err || order==nil")
		fmt.Println(info.err, userPay.err, eventPay.err, users.err, order == nil)
		return
	}
	history, pending := f.PreparePayments(userPay.payments, eventPay.payments, order)
	f.Respond(w, "ok",
		"room_name", info.roomName,
		"room_desc", info.desc,
		"history_payments", history,
		"pending_payments", pending,
		"saldo", info.saldo,
		"am_admin", info.isAdmin,
		"users", users.users,
		"users_role", users.usersIsAdmin) // map[username]isAdmin
}

func getMyRooms(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	uname := params["username"]
	retch := make(chan roomsInfo, 1)
	defer close(retch)

	sqlGetMyRooms(retch, uname)
	myrooms := <-retch

	if myrooms.err {
		f.RIE(w)
		return
	}

	f.Respond(w, "ok",
		"my_rooms", myrooms.rooms)
}

var (
	GetRoomData = http.HandlerFunc(getRoomData)
	GetMyRooms  = http.HandlerFunc(getMyRooms)
)
