package scm

import "database/sql"

func (c *ChansToMngr) ConnectCTM(m *Mngr) {
	/*
		init? (or main):
		    import sqlconnmanager as scm
			import package_x as x
		    	(package_x:
					import sqlconnmanager as scm
					var Ctm = &ChansToMngr{}
				)
			var m = InitMngr(x, c, a, s)
			x.Ctm.ConnectCTM(m)
	*/
	c.Lock.Lock()
	c.StmtChan = m.stmtChan
	c.Lock.Unlock()
}

func (c *ChansToMngr) Exec(sqlTxt string, args []interface{}) (sql.Result, error) {
	ch := make(chan QueryRet, 1)
	c.StmtChan <- InStmt{
		SqlText:    sqlTxt,
		StmtArgs:   args,
		ActionType: EXEC,
		ResChan:    ch,
	}
	qret := <-ch
	return qret.Result, qret.Err
}

func (c *ChansToMngr) Query(sqlTxt string, args []interface{}) (*sql.Rows, error) {
	ch := make(chan QueryRet, 1)
	c.StmtChan <- InStmt{
		SqlText:    sqlTxt,
		StmtArgs:   args,
		ActionType: QUERY,
		ResChan:    ch,
	}
	qret := <-ch
	return qret.Rows, qret.Err
}

func (c *ChansToMngr) QueryRow(sqlTxt string, args []interface{}) (*sql.Row, error) {
	ch := make(chan QueryRet, 1)
	c.StmtChan <- InStmt{
		SqlText:    sqlTxt,
		StmtArgs:   args,
		ActionType: QUERYROW,
		ResChan:    ch,
	}
	qret := <-ch
	return qret.Row, qret.Err
}
