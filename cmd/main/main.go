package main

import (
	"fmt"
	"net/http"
	"sync"

	cmt "github.com/kacpekwasny/commontools"
	apir "github.com/kacpekwasny/payments-backend/internal/api_router"
	"github.com/kacpekwasny/payments-backend/internal/router"
	scm "github.com/kacpekwasny/payments-backend/pkg/sql_conn_manager"
)

func main() {
	fmt.Println("running main...")
	m := scm.InitMngr(CONFIG_MAP)
	connectWrappersRouterCtm(m)
	configWrappers(CONFIG_MAP)

	// rtr -> internal router
	// wrpr -> external router that loggs requests before passing them to 'rtr'
	rtr, wrpr := router.NewRouter()
	SetRoutes(rtr)

	apir, apiwrpr := apir.NewRouter()
	SetRoutesApiRouter(apir)

	cmt.Pc(m.Connect())
	m.WatchQueries(true)

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		listen := CONFIG_MAP["listen address"]
		fmt.Printf("ListenAndServe external... on %#v\n", listen)
		err := http.ListenAndServe(listen, wrpr)
		fmt.Println("end ListenAndServe external err:", err)
	}()
	go func() {
		defer wg.Done()
		apiListen := CONFIG_MAP["listen address api"]
		fmt.Printf("ListenAndServe internal API... on %#v\n", apiListen)
		err := http.ListenAndServe(apiListen, apiwrpr)
		fmt.Println("end ListenAndServe internal API err:", err)
	}()
	wg.Wait()
}
