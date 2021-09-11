package scm

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/kacpekwasny/payments-backend/internal/funcs"
)

func InitMngr(m map[string]string) *Mngr {
	mngr := &Mngr{
		conf: SQLConnConf{
			username:     m["username"],
			password:     m["password"],
			address:      m["address"],
			databaseName: m["databaseName"],
			port:         3306,
		},
		connectionLive: false,
		stmtChan:       make(chan InStmt),
		Lock:           &sync.Mutex{},
	}
	return mngr
}

// Connection funcs //
func (m *Mngr) Connect() error {
	db, err := sql.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?parseTime=True",
			m.conf.username, m.conf.password, m.conf.address, m.conf.port, m.conf.databaseName))
	if err != nil {
		return err
	}
	m.db = db
	return nil
}

func (m *Mngr) ControlConnection(turnOn bool) {
	m.pingOnMtx.Lock()
	defer func() {
		m.pingOn = turnOn
		m.pingOnMtx.Unlock()
	}()
	// Make sure only one go routine is Controling connection
	if turnOn && !m.pingOn {
		go func() {
			// set test ping as OFF when returning
			defer func() {
				m.pingOnMtx.Lock()
				m.pingOn = false
				m.pingOnMtx.Unlock()
			}()
			m.Log1(" ControlConnectionDB() - START")
			pch := make(chan error, 1)
			mu := sync.Mutex{}
			counter := 0
			for {
				counter = counter % 10
				go func() {
					pch <- m.db.Ping()
				}()

				select {
				case err := <-pch:
					mu.Lock()
					if err != nil {
						m.Log1("Lost connection")
						m.connectionLive = false
					} else if !m.connectionLive {
						m.Log1("Regained connection")
						m.connectionLive = true
					}
					mu.Unlock()
				case <-time.After(m.pingTimeout):
					mu.Lock()
					if m.connectionLive {
						m.Log1("Lost connection")
					} else {
						// dont spam "No connection to database"
						// print "db disconnected" every couple tries
						if counter == 0 {
							m.Log1("No connection to database")
						}
					}
					m.connectionLive = false
					mu.Unlock()
				}
				time.Sleep(m.pingInterval)
				counter++
			}
		}()
	}
}

func (m *Mngr) GetConnLive() bool {
	return m.connectionLive
}

// Listen & execute queries //
func (m *Mngr) WatchQueries(turnOn bool) {
	m.watchQueriesOnMtx.Lock()
	defer func() {
		m.watchQueriesOn = turnOn
		m.watchQueriesOnMtx.Unlock()
	}()
	if turnOn && !m.watchQueriesOn {
		go func() {
			for {

				select {
				case st := <-m.stmtChan:
					qret := QueryRet{}
					switch st.ActionType {
					case 0:
						qret.Result, qret.Err = m.db.Exec(st.SqlText, st.StmtArgs...)
					case 1:
						qret.Rows, qret.Err = m.db.Query(st.SqlText, st.StmtArgs...)
					case 2:
						qret.Row = m.db.QueryRow(st.SqlText, st.StmtArgs...)
					}
					st.ResChan <- qret
					close(st.ResChan)
				case <-time.After(time.Millisecond * 100):
					// So when the m.watchQueriesOn changes select wont be stuck waiting
					// for stmtChan
				}

				m.watchQueriesOnMtx.Lock()
				if !m.watchQueriesOn {
					m.watchQueriesOnMtx.Unlock()
					return
				}
			}
		}()
	}
}

// Log funcs //
func (m *Mngr) Log1(str string, values ...interface{}) {
	if m.LOG_LEVEL >= 1 {
		funcs.Log("Mngr: "+str, values...)
	}
}

func (m *Mngr) Log2(str string, values ...interface{}) {
	if m.LOG_LEVEL >= 2 {
		funcs.Log("Mngr: "+str, values...)
	}
}

func (m *Mngr) Log3(str string, values ...interface{}) {
	if m.LOG_LEVEL >= 3 {
		funcs.Log("Mngr: "+str, values...)
	}
}
