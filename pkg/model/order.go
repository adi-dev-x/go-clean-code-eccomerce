package model

import "net/url"

type ListAllOrders struct {
	Name     string  `json:"name"`
	Unit     int     `json:"unit"`
	Status   string  `json:"status"`
	Returned bool    `json:"returned"`
	Amount   float64 `json:"amount"`
	Pid      string  `json:"pid"`
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
