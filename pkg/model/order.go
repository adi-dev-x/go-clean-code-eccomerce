package model

import (
	"net/url"
	"time"
)

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
	ListDate string  `json:"check"`
}
type SalesReport struct {
	Type string `json:"type"`
	Usid string
	From string `json:"from"`
	To   string `json:"to"`
}
type Salesfact struct {
	Revenue       float64 `json:"revenue"`
	TotalDiscount float64 `json:"total_discount"`
	TotalSales    float64 `json:"total_sales"`
	TotalOrders   int     `json:"total_orders"`
	Date          string
}

func (s *SalesReport) Valid() (err url.Values) {
	err = url.Values{}
	if !(s.Type == "Weekly" || s.Type == "Daily" || s.Type == "Yearly" || s.Type == "Monthly" || s.Type == "Custom") {
		err.Add("Wrong format of Type", " Type should be in Weekly Daily Yearly Monthly Custom")
	}
	if s.Type == "Custom" {
		const dateFormat = "2006-01-02"
		fromDate, fromErr := time.Parse(dateFormat, s.From)
		toDate, toErr := time.Parse(dateFormat, s.To)

		if fromErr != nil {
			err.Add("From", "From date should be in the format YYYY-MM-DD")
		}
		if toErr != nil {
			err.Add("To", "To date should be in the format YYYY-MM-DD")
		}

		if fromErr == nil && toErr == nil && fromDate.After(toDate) {
			err.Add("Date Range", "From date should not be greater than To date")
		}
	}

	return err
}
