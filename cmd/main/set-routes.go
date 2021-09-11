package main

import (
	"github.com/gorilla/mux"
	"github.com/kacpekwasny/payments-backend/internal/router"
	"github.com/kacpekwasny/payments-backend/internal/wrappers"
	scm "github.com/kacpekwasny/payments-backend/pkg/sql_conn_manager"
)

func SetRoutes(rtr *mux.Router) {
	const rut = "{room_link}/{username}/{token}"
	// get-routes
	rtr.Handle("/get/room-data/"+rut, wrappers.AA((router.GetRoomData))).Methods("GET")
	rtr.Handle("/get/room-payments/"+rut, wrappers.AA(router.GetRoomPayments)).Methods("GET")

	// post-routes
	rtr.Handle("/post/new-payment/"+rut, wrappers.AA(router.PostNewPayment)).Methods("POST")
	rtr.Handle("/post/add-user-room/"+rut+"/{new_user}", wrappers.AA(router.PostAddUserRoom)).Methods("POST")
	rtr.Handle("/post/remove-user-room/"+rut+"/{old_user}", wrappers.AA(router.PostRemoveUserRoom)).Methods("POST")

	// update-routes
	rtr.Handle("/update/accept-payment/"+rut+"/{payment_id}", wrappers.AA(router.UpdateAcceptPayment)).Methods("UPDATE")
	rtr.Handle("/update/unaccept-payment/"+rut+"/{payment_id}", wrappers.AA(router.UpdateUnacceptPayment)).Methods("UPDATE")
	rtr.Handle("/update/change-user-role/"+rut+"/{user_for_change}/{new_role}", wrappers.AA(router.UpdateUnacceptPayment)).Methods("UPDATE")
}

func connectWrappersRouterCtm(m *scm.Mngr) {
	wrappers.Ctm.ConnectCTM(m)
	router.Ctm.ConnectCTM(m)
}
