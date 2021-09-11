package scm

import (
	"database/sql"
	"sync"
	"time"
)

type SQLConnConf struct {
	username     string
	password     string
	address      string
	databaseName string
	tableName    string
	port         int
}

type Mngr struct {
	conf SQLConnConf
	Lock *sync.Mutex

	db                *sql.DB
	watchQueriesOn    bool
	watchQueriesOnMtx *sync.Mutex

	connectionLive bool
	pingOnMtx      *sync.Mutex
	pingOn         bool
	pingTimeout    time.Duration // time waiting for ping to be returned
	pingInterval   time.Duration // how often pings are sent

	// In statements
	stmtChan chan InStmt

	LOG_LEVEL int
}

type InStmt struct {
	SqlText    string
	StmtArgs   []interface{}
	ActionType int // 0: Stmt.Exec, 1: Stmt.Query, 2: Stmt.QueryRow
	ResChan    chan QueryRet
}

type QueryRet struct {
	// Struct holding output of stmt.Exec() or stmt.Query() or stmt.QueryRow()
	Result sql.Result
	Rows   *sql.Rows
	Row    *sql.Row
	Err    error
}

type ChansToMngr struct {
	// this struct will hold chanels that allow communication
	// to a Mngr which will execute queries to an SQL DB
	StmtChan chan InStmt
	Lock     *sync.Mutex
}
