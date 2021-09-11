package jstr

import "time"

type Rooms struct {
	Link string `json:"link"`
	Name string `json:"name"`
	Desc string `json:"desc"`
}

type Users struct {
	Username string `json:"username"`
}

type PaymentFromTo struct {
	ID           int       `json:"id"`
	Title        string    `json:"title"`
	Value        float64   `json:"value"`
	Created      time.Time `json:"created"`
	FromAccepted bool      `json:"from_accepted"`
	ToAccepted   bool      `json:"to_accepted"`
	FromUsername string    `json:"from_username"`
	ToUsername   string    `json:"to_username"`
}
