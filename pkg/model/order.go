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
	Data      []ResultsAdminsalesReport
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
	Amount   float64 `json:"Product_price"`
	Date     string  `json:"date"`
	Oid      string  `json:"oid"`
	VName    string  `json:"vendor_name"`
	Discount float64 `json:"discount"`
	Cmt      float64 `json:"Coupon_Amount"`
	Code     string  `json:"Coupon_code"`
	Wmt      float64 `json:"wallet_amount"`
}
type ResultsAdminsalesReport struct {
	Name     string  `json:"name"`
	Unit     int     `json:"unit"`
	Amount   float64 `json:"Total_Amount"`
	Date     string  `json:"date"`
	Oid      string  `json:"oid"`
	VName    string  `json:"vendor_name"`
	Discount float64 `json:"discount"`
	Cmt      float64 `json:"Coupon_Amount"`
	Code     string  `json:"code"`
}
type ResultsVendorsales struct {
	Name   string  `json:"name"`
	Unit   int     `json:"unit"`
	Amount float64 `json:"Total_amount"`
	Date   string  `json:"date"`
	Oid    string  `json:"order_id"`

	Discount float64 `json:"discount"`
	Cmt      float64 `json:"coupon_Amount"`
	Code     string  `json:"code"`
}
type ListOrdersVendor struct {
	Name       string  `json:"name"`
	Unit       int     `json:"unit"`
	Status     string  `json:"status"`
	Returned   bool    `json:"returned"`
	Amount     float64 `json:"Product_price"`
	Pid        string  `json:"pid"`
	Date       string  `json:"date"`
	User       string  `json:"user"`
	Add        string  `json:"user_ad"`
	ListDate   string  `json:"checks"`
	Oid        string  `json:"oid"`
	Discount   float64 `json:"discount"`
	CouponAmt  float64 `json:"coupon_amount"`
	CouponCode string  `json:"coupon_code"`
	WalletAmt  float64 `json:"wallet_amount"`
}
type ListOrdersAdmin struct {
	Name       string  `json:"name"`
	Unit       int     `json:"unit"`
	Status     string  `json:"status"`
	Returned   bool    `json:"returned"`
	Amount     float64 `json:"Product_price"`
	Pid        string  `json:"pid"`
	Date       string  `json:"date"`
	User       string  `json:"user"`
	Add        string  `json:"user_ad"`
	ListDate   string  `json:"checks"`
	Oid        string  `json:"oid"`
	VName      string  `json:"vendorname"`
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

type UpdateOrderAdmin struct {
	Oid            string `json:"oid"`
	Delivery_date  string `json:"delivery_date"`
	Payment_status string `json:"payment_status"`
	Delivery_Stat  string `json:"delivered"`
}

func (s *UpdateOrderAdmin) Valid() (err url.Values) {
	err = url.Values{}
	if s.Oid == "" {
		err.Add("order id", "it should not be nil")

	}
	if s.Delivery_date != "" {
		const dateFormat = "2006-01-02"
		Date, DErr := time.Parse(dateFormat, s.Delivery_date)
		if DErr != nil {
			err.Add("From", "From date should be in the format YYYY-MM-DD")
		} else {
			currentDate := time.Now().Truncate(24 * time.Hour)

			// Compare delivery date with the current date
			if Date.Before(currentDate) {
				err.Add("Delivery_date", "Delivery date should not be before the current date")
			}
		}

	}
	if !(s.Delivery_Stat == "Delivered" || s.Delivery_Stat == "Not Delivered" || s.Delivery_Stat == "") {
		err.Add("Wrong format of Delivery status", " Delivery status should be in Delivered, Not Delivered or Blank")
	}
	if !(s.Payment_status == "Pending" || s.Payment_status == "Completed" || s.Payment_status == "Failed" || s.Payment_status == "") {
		err.Add("Wrong format of Payment status", " Payment status should be in Pending Completed Failed or Blank")
	}

	return err

}
