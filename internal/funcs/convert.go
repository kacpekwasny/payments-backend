package funcs

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	jstr "github.com/kacpekwasny/payments-backend/internal/json_structs"
)

func PrepareEventPayments(rows *sql.Rows) (map[string]jstr.Payment, error) {
	epM := map[string]jstr.Payment{}

	ep := jstr.Payment{} // EventPayment
	ep.Data.ID = ""
	ep.Data.Ammounts = map[string]float64{}
	ep.Data.Accepted = []string{}
	ep.Data.Waiting = []string{}

	ept := jstr.Payment{} // EvenyPayment temorary
	rowsCount := 0
	for rows.Next() {
		rowsCount++
		t := time.Time{}
		uname := ""
		accepted := false
		input := 0.0
		//	SELECT ep.public_id, ep.title, ep.created,
		//		(SELECT username FROM payments.users AS u WHERE u.id=epd.users_id) as username,
		//	epd.user_has_accepted, epd.ammount, ep.all_users_accepted
		err := rows.Scan(&ept.Data.ID, &ept.Data.Title, &t, &uname, &accepted, &input, &ept.Data.AllUsersAccepted)
		if err != nil {
			return nil, err
		}

		if ep.Data.ID != "" && ept.Data.ID != ep.Data.ID {
			// from this loop and forward we are parsing another event payment

			// if all users accepted, set type as history
			if len(ep.Data.Waiting) == 0 {
				ep.HisPend = jstr.HISTORY
			} else { // else set type as pending
				ep.HisPend = jstr.PENDING
			}
			cp := ep
			epM[cp.Data.ID] = cp

			// parsing new Payment
			ep = jstr.Payment{} // EventPayment
			ep.Data.Ammounts = map[string]float64{}
			ep.Data.Accepted = []string{}
			ep.Data.Waiting = []string{}
		}
		ep.Type = jstr.EVENT_PAYMENT
		ep.Data.ID = ept.Data.ID
		ep.Data.Title = ept.Data.Title
		ep.Data.Created = t.Format("02.01.2006")
		ep.Data.AllUsersAccepted = ept.Data.AllUsersAccepted
		// set ammount
		ep.Data.Ammounts[uname] = input

		// set accepted
		if accepted {
			ep.Data.Accepted = append(ep.Data.Accepted, uname)
		} else {
			ep.Data.Waiting = append(ep.Data.Waiting, uname)
		}
	}
	if len(ep.Data.Waiting) == 0 {
		ep.HisPend = jstr.HISTORY
	} else { // else set type as pending
		ep.HisPend = jstr.PENDING
	}
	if rowsCount > 0 {
		epM[ep.Data.ID] = ep
	}
	return epM, nil
}

func PrepareUserPayments(rows *sql.Rows) (map[string]jstr.Payment, error) {
	upM := map[string]jstr.Payment{}
	for rows.Next() {
		t := time.Time{}
		fromAccepted := false
		toAccepted := false
		up := jstr.Payment{
			Type: jstr.USER_PAYMENT,
		}
		up.Data.Accepted = []string{}
		up.Data.Waiting = []string{}
		err := rows.Scan(&up.Data.ID, &up.Data.Title, &up.Data.Value, &t, &fromAccepted, &toAccepted,
			&up.Data.AllUsersAccepted, &up.Data.FromUsername, &up.Data.ToUsername)
		if err != nil {
			return nil, err
		}
		if fromAccepted {
			up.Data.Accepted = append(up.Data.Accepted, up.Data.FromUsername)
		} else {
			up.Data.Waiting = append(up.Data.Waiting, up.Data.FromUsername)
		}
		if toAccepted {
			up.Data.Accepted = append(up.Data.Accepted, up.Data.ToUsername)
		} else {
			up.Data.Waiting = append(up.Data.Waiting, up.Data.ToUsername)
		}
		if fromAccepted && toAccepted {
			up.HisPend = jstr.HISTORY
		} else {
			up.HisPend = jstr.PENDING
		}
		up.Data.Created = t.Format("02.01.2006")
		upM[up.Data.ID] = up
	}
	return upM, nil
}

func PreparePaymentsOrder(rows *sql.Rows) ([]string, error) {
	order := []string{}
	for rows.Next() {
		public_id := ""
		created := time.Time{}
		err := rows.Scan(&public_id, &created)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		order = append(order, public_id)
	}
	return order, nil
}

// return history, pending, error
func PreparePayments(userP, eventP map[string]jstr.Payment, order []string) ([]jstr.Payment, []jstr.Payment) {
	history := []jstr.Payment{}
	pending := []jstr.Payment{}
	for _, id := range order {
		up, ok := userP[id]
		if ok {
			if up.Data.AllUsersAccepted {
				history = append(history, up)
				continue
			}
			pending = append(pending, up)
			continue
		}
		ep, ok := eventP[id]
		if ok {
			if ep.Data.AllUsersAccepted {
				history = append(history, ep)
				continue
			}
			pending = append(pending, ep)
			continue
		}
	}
	return history, pending
}

func PreparePaymentFromRequest(r *http.Request) (jstr.Payment, error) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return jstr.Payment{}, err
	}

	PPJsonBytes(b)

	var paymentTmp jstr.Payment
	err = json.Unmarshal(b, &paymentTmp)

	return paymentTmp, err
}

func PrepareNameDescFromRequest(r *http.Request) (string, string, error) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return "", "", err
	}

	PPJsonBytes(b)

	var mp = map[string]string{}
	err = json.Unmarshal(b, &mp)
	if err != nil {
		return "", "", err
	}

	roomName, _ := mp["roomName"]
	roomDesc, _ := mp["roomDesc"]
	return roomName, roomDesc, nil
}
