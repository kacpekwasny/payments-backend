package wrappers

import (
	"net/http"
	"sync"
	"time"

	scm "github.com/kacpekwasny/payments-backend/pkg/sql_conn_manager"
)

var (
	Ctm = scm.ChansToMngr{
		Lock: &sync.Mutex{},
	}
	Config = configStruct{
		Lock: &sync.Mutex{},
	}
	client = &http.Client{
		Timeout: time.Second,
	}
)

type configStruct struct {
	AuthApiBaseUrl string
	Lock           *sync.Mutex
	AuthON         bool
}
