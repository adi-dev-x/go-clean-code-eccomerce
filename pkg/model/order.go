package model

import "net/url"

type ListAllOrders struct {
	Name     string  `json:"name"`
	Unit     int     `json:"unit"`
	Status   string  `json:"status"`
	Returned bool    `json:"returned"`
	Amount   float64 `json:"amount"`
	Pid      string  `json:"pid"`
	Date     string  `json:"date"`
}
type ListAllOrdersUsers struct {
	Name     string  `json:"name"`
	Unit     int     `json:"unit"`
	Status   string  `json:"status"`
	Returned bool    `json:"returned"`
	Amount   float64 `json:"amount"`
	Pid      string  `json:"pid"`
	Date     string  `json:"date"`
	Oid      string  `json:"oid"`
}
type ListAllOrdersCheck struct {
	Name     string  `json:"name"`
	Unit     int     `json:"unit"`
	Status   string  `json:"status"`
	Returned bool    `json:"returned"`
	Amount   float64 `json:"amount"`
	Pid      string  `json:"pid"`
	Date     string  `json:"date"`
	Usid     string  `json:"user_id"`
	Vid      string  `json:"vid"`
	Usmail   string  `json:"usmail"`
}

type ReturnOrderPost struct {
	Oid string `json:"oid"`
}

func (r *ReturnOrderPost) Valid() (err url.Values) {
	err = url.Values{}
	if r.Oid == "" {
		err.Add("item id ", "id is not present")

	}

	return err

}

type ListOrdersVendor struct {
	Name     string  `json:"name"`
	Unit     int     `json:"unit"`
	Status   string  `json:"status"`
	Returned bool    `json:"returned"`
	Amount   float64 `json:"amount"`
	Pid      string  `json:"pid"`
	Date     string  `json:"date"`
	User     string  `json:"user"`
	Add      string  `json:"user_ad"`
}
