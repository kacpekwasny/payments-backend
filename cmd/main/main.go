package main

import (
	"fmt"
	"net/http"

	"github.com/kacpekwasny/payments-react/internal/router"
	scm "github.com/kacpekwasny/payments-react/pkg/sql_conn_manager"
)

func main() {
	fmt.Println("running main...")
	m := scm.InitMngr(CONFIG_MAP)
	connectWrappersRouterCtm(m)

	rtr, wrpr := router.NewRouter()
	// rtr -> internal router
	// wrpr -> external router that loggs requests before passing them to 'rtr'
	SetRoutes(rtr)
	fmt.Println("ListenAndServe...")
	http.ListenAndServe(CONFIG_MAP["listen port"], wrpr)
	fmt.Println("end ListenAndServe")
}
