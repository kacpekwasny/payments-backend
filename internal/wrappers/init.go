package wrappers

import (
	"net/http"
	"time"

	scm "github.com/kacpekwasny/payments-react/pkg/sql_conn_manager"
)

var (
	Ctm    = scm.ChansToMngr{}
	config = configStruct{}
	client = &http.Client{
		Timeout: time.Second,
	}
)

type configStruct struct {
	AuthApiBaseUrl string
}
