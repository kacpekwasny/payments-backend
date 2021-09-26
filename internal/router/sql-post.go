package router

import (
	"fmt"

	"github.com/kacpekwasny/payments-backend/internal/json_structs"
)

// return 'ok'
func insertUserPayment(roomsLink, createdBy string, pay json_structs.Payment) bool {

	_, err := Ctm.Exec(`CALL InsertUserPayment(?, ?, ?, ?, ?, ?)`,
		[]interface{}{roomsLink, pay.Data.FromUsername, pay.Data.ToUsername,
			pay.Data.Value, pay.Data.Title, createdBy})

	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
	//rowsNum, _ := result.RowsAffected()
	//if rowsNum != 1 {
	//	fmt.Println("rowsNum: ", rowsNum, "roomsLink: ", roomsLink, "createdBy: ", createdBy, "pay:", pay)
	//	return false
	//}
	//return true
}

func insertEventPayment(roomsLink, createdBy string, pay json_structs.Payment) bool {
	row, err := Ctm.QueryRow(`CALL InsertEventPayment(?, ?, ?)`,
		[]interface{}{roomsLink, pay.Data.Title, createdBy})

	if err != nil {
		fmt.Println("Ctm.QueryRow(...)", err)
		return false
	}
	publicID := ""
	err = row.Scan(&publicID)
	if err != nil {
		fmt.Println("row.Scan(publicID)", err)
		return false
	}
	for uname, amm := range pay.Data.Ammounts {
		go Ctm.Exec(`CALL InsertEventPaymentPUD(?, ?, ?)`,
			[]interface{}{publicID, uname, amm})
	}
	return true
}
