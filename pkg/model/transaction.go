package model

type UserTransactions struct {
	Id     string  `json:"id"`
	Amount float64 `json:"amount"`
	Type   string  `json:"type"`
}
