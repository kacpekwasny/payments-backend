package api_router

import (
	"sync"

	scm "github.com/kacpekwasny/payments-backend/pkg/sql_conn_manager"
)

var (
	Ctm = scm.ChansToMngr{
		Lock: &sync.Mutex{},
	}
)
