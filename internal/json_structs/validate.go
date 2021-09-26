package json_structs

/*
functions to check if user sent data is valid to be saved in databases
these function concern Payments
*/

func UserPaymentIsValid(p Payment) bool {
	d := p.Data

	if p.Type != USER_PAYMENT ||
		len(d.FromUsername) > 40 ||
		len(d.ToUsername) > 40 ||
		len(d.Title) > 100 ||
		d.ToUsername == d.FromUsername ||
		d.Value <= 0 {
		return false
	}

	return true
}

func EventPaymentIsValid(p Payment) bool {
	d := p.Data

	if p.Type != EVENT_PAYMENT ||
		len(d.Title) > 100 ||
		len(d.Ammounts) < 2 {
		return false
	}

	for uname, amm := range d.Ammounts {
		if len(uname) < 1 ||
			len(uname) > 40 ||
			amm < 0 {
			return false
		}
	}

	return true
}
