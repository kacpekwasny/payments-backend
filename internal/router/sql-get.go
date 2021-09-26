package router

import (
	"database/sql"
	"fmt"
	"time"

	f "github.com/kacpekwasny/payments-backend/internal/funcs"
	jstr "github.com/kacpekwasny/payments-backend/internal/json_structs"
)

type sqlGetInfo struct {
	roomName string
	desc     string
	saldo    float64
	isAdmin  bool
	err      bool
}

func sqlGetRoomInfo(retch chan sqlGetInfo, link, uname string) {

	row, err := Ctm.QueryRow(`
		SELECT r.`+"`name`"+`, r.`+"`desc`"+`, rhu.saldo, rhu.is_admin
		FROM payments.rooms as r
		INNER JOIN payments.rooms_have_users as rhu
			ON r.id=rhu.rooms_id
		WHERE
			r.link=?
		AND
			rhu.users_id=(SELECT id FROM payments.users WHERE username=?)`,
		[]interface{}{link, uname})

	if err != nil {
		retch <- sqlGetInfo{err: true}
		fmt.Println("getRoomData, qret.Err:", err)
		return
	}
	var (
		roomName, desc string
		saldo          float64
		isAdmin        bool
	)
	err = row.Scan(&roomName, &desc, &saldo, &isAdmin)
	if err != nil {
		retch <- sqlGetInfo{err: true}
		fmt.Println("row.Scan(&roomName, &desc, &saldo, &isAdmin)", err, uname, link)
		return
	}
	retch <- sqlGetInfo{roomName, desc, saldo, isAdmin, false}
}

// Get Payments from sql
type sqlGetPay struct {
	payments map[string]jstr.Payment
	err      bool
}

func sqlGetUserPayments(retch chan sqlGetPay, link string) {

	rows, err := Ctm.Query(`SELECT public_id, title, ammount, created, from_accepted, to_accepted, all_users_accepted,
		(SELECT username FROM users WHERE id=up.from_id) as 'from',
		(SELECT username FROM users WHERE id=up.to_id) as 'to'
		FROM user_payments as up
		WHERE up.rooms_link=?
		ORDER BY up.created ASC`,
		[]interface{}{link})

	if err != nil {
		if err == sql.ErrNoRows {
			retch <- sqlGetPay{nil, false}
			return
		}
		retch <- sqlGetPay{nil, true}
		panic(err)
	}

	up, err := f.PrepareUserPayments(rows)
	if err != nil {
		retch <- sqlGetPay{nil, true}
		panic(err)
	}
	retch <- sqlGetPay{up, false}
}

func sqlGetEventPayments(retch chan sqlGetPay, link string) {

	rows, err := Ctm.Query(`SELECT ep.public_id, ep.title, ep.created, 
		(SELECT username FROM payments.users AS u WHERE u.id=epd.users_id) as `+"`username`"+`,
		epd.user_has_accepted, epd.ammount, ep.all_users_accepted
		FROM payments.event_payments as ep
		INNER JOIN payments.event_payments_per_user_detail as epd
			ON ep.id=epd.event_payments_id
		WHERE ep.rooms_link=?
		ORDER BY ep.id`,
		[]interface{}{link})

	if err != nil {
		if err == sql.ErrNoRows {
			retch <- sqlGetPay{nil, true}
			return
		}
		retch <- sqlGetPay{nil, false}
		panic(err)
	}

	ep, err := f.PrepareEventPayments(rows)
	if err != nil {
		retch <- sqlGetPay{nil, true}
		panic(err)
	}
	retch <- sqlGetPay{ep, false}
}

func sqlGetPaymentsOrder(retch chan []string, link string) {
	rows, err := Ctm.Query(`
		(SELECT public_id, created
		FROM user_payments WHERE rooms_link=?)
			UNION
		(SELECT public_id, created
		FROM event_payments WHERE rooms_link=?)
			ORDER BY created DESC`,
		[]interface{}{link, link})

	if err != nil {
		retch <- nil
		return
	}
	order, err := f.PreparePaymentsOrder(rows)
	if err != nil {
		fmt.Println(err)
	}
	retch <- order
}

type roomsInfo struct {
	err   bool
	rooms []room
}

type room struct {
	Name    string    `json:"name"`
	Link    string    `json:"link"`
	Saldo   float64   `json:"saldo"`
	IsAdmin bool      `json:"is_admin"`
	Joined  time.Time `json:"joined"`
}

func sqlGetMyRooms(retch chan roomsInfo, uname string) {
	rows, err := Ctm.Query(`
		SELECT r.name, r.link, rhu.saldo, rhu.is_admin, rhu.joined_room
		FROM rooms_have_users as rhu
		INNER JOIN rooms as r
			ON
		rhu.rooms_id = r.id 
		INNER JOIN users as u
			ON
		rhu.users_id=u.id
		WHERE u.username=?
		ORDER BY rhu.joined_room DESC`,
		[]interface{}{uname})
	if err != nil {
		fmt.Println(err)
		retch <- roomsInfo{err: true}
		return
	}

	roomls := []room{}
	for rows.Next() {
		r := room{}
		err = rows.Scan(&r.Name, &r.Link, &r.Saldo, &r.IsAdmin, &r.Joined)
		if err != nil {
			fmt.Println(err)
			retch <- roomsInfo{err: true}
			return
		}
		roomls = append(roomls, r)
	}
	retch <- roomsInfo{
		err:   false,
		rooms: roomls,
	}
}

type usersInfo struct {
	err          bool
	users        []string
	usersIsAdmin map[string]bool
}

func sqlGetRoomUsers(retch chan usersInfo, roomLink string) {
	rows, err := Ctm.Query(`
		SELECT
			(SELECT username FROM payments.users WHERE id=rhu.users_id) as "username",
			is_admin
		FROM payments.rooms_have_users as rhu
		WHERE rooms_link=?`,
		[]interface{}{roomLink})
	if err != nil {
		fmt.Println(err)
		retch <- usersInfo{err: true}
		return
	}
	users := []string{}
	usersIsAdmin := map[string]bool{}
	for rows.Next() {
		username := ""
		isAdmin := false
		err = rows.Scan(&username, &isAdmin)
		if err != nil {
			fmt.Println(err)
			retch <- usersInfo{err: true}
			return
		}
		users = append(users, username)
		usersIsAdmin[username] = isAdmin
	}
	retch <- usersInfo{
		err:          false,
		users:        users,
		usersIsAdmin: usersIsAdmin,
	}
}
