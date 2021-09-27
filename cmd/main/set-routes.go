package main

import (
	"strconv"

	"github.com/gorilla/mux"
	cmt "github.com/kacpekwasny/commontools"
	apir "github.com/kacpekwasny/payments-backend/internal/api_router"
	"github.com/kacpekwasny/payments-backend/internal/router"
	"github.com/kacpekwasny/payments-backend/internal/wrappers"
	scm "github.com/kacpekwasny/payments-backend/pkg/sql_conn_manager"
)

func SetRoutes(rtr *mux.Router) {
	const rut = "{room_link}/{username}/{token}"
	// get-routes
	rtr.Handle("/get/room-data/"+rut, wrappers.AA(router.GetRoomData)).Methods("GET")
	rtr.Handle("/get/my-rooms/{username}/{token}", wrappers.AA(router.GetMyRooms)).Methods("GET")

	// post-routes
	// rtr.Handle("/post/new-room/{username}/{token}", wrappers.UserIsAuthenticated(router.PostNewRoom)).Methods("POST")
	rtr.Handle("/post/new-room/{username}/{token}", wrappers.AA(router.PostNewRoom)).Methods("POST") // <- temporary
	rtr.Handle("/post/new-payment/"+rut, wrappers.AA(router.PostNewPayment)).Methods("POST")
	rtr.Handle("/post/add-user-room/"+rut+"/{new_user}", wrappers.AA(router.PostAddUserRoom)).Methods("POST")
	rtr.Handle("/post/remove-user-room/"+rut+"/{old_user}", wrappers.AA(router.PostRemoveUserRoom)).Methods("POST")
	rtr.Handle("/post/leave-room/"+rut, wrappers.AA(router.LeaveRoom)).Methods("POST")

	// update-routes
	rtr.Handle("/update/set-accept-payment/"+rut+"/{payment_id}/{accept}", wrappers.AA(router.UpdateSetAcceptPayment)).Methods("PUT")
	rtr.Handle("/update/change-user-role/"+rut+"/{user_for_change}/{new_role}", wrappers.AA(router.UpdateChangeUserRole)).Methods("PUT")
	rtr.Handle("/update/room-name-desc/"+rut, wrappers.AA(router.UpdateRoomNameAndDescription)).Methods("PUT")
	// update Room Name
	// update Room Description
}

func SetRoutesApiRouter(rtr *mux.Router) {
	rtr.Handle("/post/change-login/{old_login}/{new_login}", apir.UpdateLogin).Methods("POST")
	rtr.Handle("/post/add-user/{username}", apir.AddUser).Methods("POST")
}

func connectWrappersRouterCtm(m *scm.Mngr) {
	wrappers.Ctm.ConnectCTM(m)
	router.Ctm.ConnectCTM(m)
	apir.Ctm.ConnectCTM(m)
}

func configWrappers(confMap map[string]string) {
	wrappers.Config.Lock.Lock()
	wrappers.Config.AuthApiBaseUrl = confMap["authApiBaseUrl"]
	ON, err := strconv.ParseBool(confMap["authON"])
	if err != nil {
		cmt.Pc(err)
	}
	wrappers.Config.AuthON = ON
	wrappers.Config.Lock.Unlock()
}
