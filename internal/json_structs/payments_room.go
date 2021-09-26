package json_structs

//type Rooms struct {
//	Link string `json:"link"`
//	Name string `json:"name"`
//	Desc string `json:"desc"`
//}

type Users struct {
	Username string `json:"username"`
}

type Payment struct {
	Type    string      `json:"type"`    // UserPayments || EventPayment
	HisPend string      `json:"HisPend"` // History || Pending
	Data    PaymentData `json:"data"`
}

type PaymentData struct {
	ID               string `json:"id"`
	paymentNum       int
	Title            string   `json:"title"`
	Value            float64  `json:"ammount"`
	Created          string   `json:"created"`
	Accepted         []string `json:"accepted"`
	Waiting          []string `json:"waiting"`
	AllUsersAccepted bool     `json:"all_users_accepted"`

	FromUsername string `json:"from_username"`
	ToUsername   string `json:"to_username"`

	Ammounts map[string]float64 `json:"ammounts"`
}

const (
	USER_PAYMENT  = "UserPayment"
	EVENT_PAYMENT = "EventPayment"

	HISTORY = "History"
	PENDING = "Pending"
)
