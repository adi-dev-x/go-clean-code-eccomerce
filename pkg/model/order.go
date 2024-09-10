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
	Discount float64 `json:"discount"`
}
type ListAllOrdersCheck struct {
	Name     string  `json:"name"`
	Unit     int     `json:"unit"`
	Status   string  `json:"status"`
	Returned string  `json:"re_cl"`
	Amount   float64 `json:"amount"`
	Pid      string  `json:"pid"`
	Date     string  `json:"date"`
	Usid     string  `json:"user_id"`
	Vid      string  `json:"vid"`
	Usmail   string  `json:"usmail"`
	Moid     string  `json:"mid"`
	Delivery string  `json:"delivered"`
}

type ReturnOrderPost struct {
	Oid      string `json:"oid"`
	MoReturn bool
}

type ReturnOrderPostForUser struct {
	Oid      string `json:"oid"`
	MoReturn bool
	Type     string `json:"type"`
}

func (r *ReturnOrderPostForUser) Valid() (err url.Values) {
	err = url.Values{}
	if r.Oid == "" {
		err.Add("item id ", "id is not present")

	}
	if !(r.Type == "Returned" || r.Type == "Cancelled") {
		err.Add("Type  ", "Give Returned or Cancelled")

	}

	return err

}
func (r *ReturnOrderPost) Valid() (err url.Values) {
	err = url.Values{}
	if r.Oid == "" {
		err.Add("item id ", "id is not present")

	}

	return err

}

type SendSalesReort struct {
	Data      []ListOrdersVendor
	FactsData Salesfact
	PdfUrl    string
	ExcelUrl  string
}
type SendSalesReortAdmin struct {
	Data      []ResultsAdminsales
	FactsData Salesfact
	PdfUrl    string
	ExcelUrl  string
}
type SendSalesReortVendorinAdmin struct {
	Data      []ResultsVendorsales
	FactsData Salesfact
	PdfUrl    string
	ExcelUrl  string
}
type ResultsAdminsales struct {
	Name     string  `json:"name"`
	Unit     int     `json:"unit"`
	Amount   float64 `json:"amount"`
	Date     string  `json:"date"`
	Oid      string  `json:"oid"`
	VName    string  `json:"vname"`
	Discount float64 `json:"discount"`
	Cmt      float64 `json:"cmt"`
	Code     string  `json:"code"`
	Wmt      float64 `json:"wmt"`
}
type ResultsVendorsales struct {
	Name   string  `json:"name"`
	Unit   int     `json:"unit"`
	Amount float64 `json:"amount"`
	Date   string  `json:"date"`
	Oid    string  `json:"oid"`

	Discount float64 `json:"discount"`
	Cmt      float64 `json:"cmt"`
	Code     string  `json:"code"`
	Wmt      float64 `json:"wmt"`
}
type ListOrdersVendor struct {
	Name       string  `json:"name"`
	Unit       int     `json:"unit"`
	Status     string  `json:"status"`
	Returned   bool    `json:"returned"`
	Amount     float64 `json:"amount"`
	Pid        string  `json:"pid"`
	Date       string  `json:"date"`
	User       string  `json:"user"`
	Add        string  `json:"user_ad"`
	ListDate   string  `json:"checks"`
	Oid        string  `json:"oid"`
	Discount   float64 `json:"discount"`
	CouponAmt  float64 `json:"cmt"`
	CouponCode string  `json:"code"`
	WalletAmt  float64 `json:"wmt"`
}
type ListOrdersAdmin struct {
	Name       string  `json:"name"`
	Unit       int     `json:"unit"`
	Status     string  `json:"status"`
	Returned   bool    `json:"returned"`
	Amount     float64 `json:"amount"`
	Pid        string  `json:"pid"`
	Date       string  `json:"date"`
	User       string  `json:"user"`
	Add        string  `json:"user_ad"`
	ListDate   string  `json:"checks"`
	Oid        string  `json:"oid"`
	VName      string  `json:"vname"`
	Discount   float64 `json:"discount"`
	CouponAmt  float64 `json:"cmt"`
	CouponCode string  `json:"code"`
	WalletAmt  float64 `json:"wmt"`
}
type SalesReport struct {
	Type string `json:"type"`
	Usid string
	From string `json:"from"`
	To   string `json:"to"`
	Vid  string `json:"vid"`
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

type ListingMainOrders struct {
	OR_id         string  `json:"id"`
	Delivery_Stat bool    `json:"delivered"`
	D_Type        string  `json:"payment_method"`
	O_status      string  `json:"status"`
	Amount        float64 `json:"payable_amount"`
	User          string  `json:"user"`
	UserAddress   string  `json:"user_address"`
	Delivery_date string  `json:"delivery_date"`
	Discount      float64 `json:"discount"`
	Cmt           float64 `json:"cmt"`
	Code          string  `json:"code"`
	Wmt           float64 `json:"wmt"`
}
type ListingUserMainOrders struct {
	OR_id         string  `json:"id"`
	Delivery_Stat bool    `json:"delivered"`
	D_Type        string  `json:"payment_method"`
	O_status      string  `json:"status"`
	Amount        float64 `json:"payable_amount"`
	Delivery_date string  `json:"delivery_date"`
	Discount      float64 `json:"discount"`
	Cmt           float64 `json:"cmt"`
	Code          string  `json:"code"`
	Wmt           float64 `json:"wmt"`
}
