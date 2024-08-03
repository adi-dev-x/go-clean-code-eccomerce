package model

import (
	"fmt"
	"net/url"
)

type Cart struct {
	Productid string `json:"product_id"`
	Unit      int    `json:"unit"`
	Userid    string `json:"user_id"`
}
type Usercartview struct {
	Productcategory string `json:"category"`
	Productname     string `json:"name"`
	Unit            int    `json:"unit"`
	Productamount   int    `json:"amount"`
	Productstatus   bool   `json:"status"`
	Vendorname      string `json:"vendorName"`
}
type UserWishview struct {
	Productcategory string `json:"category"`
	Productname     string `json:"name"`

	Productamount int    `json:"amount"`
	Productstatus bool   `json:"status"`
	Vendorname    string `json:"vendorName"`
}

type Wishlist struct {
	Productid string `json:"product_id"`

	Userid string `json:"user_id"`
}
type Coupon struct {
	Code      string `json:"code"`
	Expiry    string `json:"expiry"`
	Minamount int    `json:"min_amount"`
	Amount    int    `json:"amount"`
}
type CouponRes struct {
	Cid         string `json:"cid"`
	Code        string `json:"code"`
	Expiry      string `json:"expiry"`
	Minamount   int    `json:"min_amount"`
	Amount      int    `json:"amount"`
	CurrentDate string `json:"current_date"`
	Is_expired  bool   `json:"is_expired"`
	Is_eligible bool   `json:"is_eligible"`
	Used        bool   `json:"used"`

	Present bool
}
type Cartresponse struct {
	Cid      string `json:"cid"`
	Pid      string `json:"pid"`
	Usid     string `json:"usid"`
	Amount   int    `json:"amount"`
	Unit     int    `json:"unit"`
	Discount int    `json:"discount"`
}
type CartresponseData struct {
	Data []Cartresponse
	Err  error
}
type FirstAddOrder struct {
	Data     CartresponseData
	TAmount  int
	CData    CouponRes
	Notvalid bool
}

type RZpayment struct {
	Id    string
	Amt   int
	Token string
}
type Order struct {
	Cartid       string `json:"cart_id"`
	Couponid     string `json:"coupon_id"`
	Type         string `json:"cod"`
	Returnstatus bool   `json:"returnstatus"`
	Aid          string `json:"aid"`
}

func (u *Order) Valid() url.Values {
	err := url.Values{}

	if u.Aid == "" {
		err.Add("Address ", "please ADD Address")
		return err
	}
	if u.Type == "" {
		err.Add("Payment Type  ", "please ADD Payment Type")
		return err
	}
	// if u.Aid == "" {
	// 	err.Add("Aid ", "no address")
	// 	return err
	// }
	return url.Values{}

}

type CheckOut struct {
	Couponid     string `json:"coupon_id"`
	Type         string `json:"cod"`
	Returnstatus bool   `json:"returnstatus"`
	Aid          string `json:"aid"`
	Wallet       bool   `json:"w_amt"`
}

func (u *CheckOut) Valid() (err url.Values, Coupon bool) {
	err = url.Values{}
	Coupon = false
	if u.Aid == "" {
		err.Add("Address ", "please ADD Address")

	}
	if !(u.Type == "ONLINE" || u.Type == "COD") {
		err.Add("Payment Type  ", "please ADD Payment Type")

	}
	// if u.Aid == "" {
	// 	err.Add("Aid ", "no address")
	// 	return err
	// }
	return err, Coupon

}

type Placeorderlist struct {
	Data string
	Err  error
}
type InsertOrder struct {
	Usid       string
	Amount     int
	Discount   int
	CouponAmt  float64
	WalletAmt  float64
	PayableAmt float64
	PayType    string
	Aid        string
	Status     string
	CouponId   string
}
type PaymentInsert struct {
	OrderId string
	Usid    string
	Amount  float64
	Status  string
	Type    string
}

func (u *CouponRes) Valid() (err url.Values) {
	err = url.Values{}

	if !u.Is_eligible {
		fmt.Println("in check 1!!!")
		err.Add("Amount ", "Total amount is less")

	}
	if u.Is_expired {
		fmt.Println("in check 2!!!")
		err.Add("Expired ", "coupon is expired")

	}
	if u.Used {
		fmt.Println("in check 3!!!")
		err.Add("Used ", "coupon is already used")

	}
	// if u.Aid == "" {
	// 	err.Add("Aid ", "no address")
	// 	return err
	// }
	return err

}
