package scm

func (c *ChansToMngr) ConnectCTM(m *Mngr) {
	/*
		init? (or main):
		    import sqlconnmanager as scm
			import package_x as x
		    	(package_x:
					import sqlconnmanager as scm
					var Ctm = &scm.ChansToMngr{}
				)
			var m = InitMngr(x, c, a, s)
			x.Ctm.ConnectCTM(m)
	*/
	c.Lock.Lock()
	c.StmtChan = m.stmtChan
	c.Lock.Unlock()
}
