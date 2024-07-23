package model

import "net/url"

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
	Cid       string `json:"cid"`
	Code      string `json:"code"`
	Expiry    string `json:"expiry"`
	Minamount int    `json:"min_amount"`
	Amount    int    `json:"amount"`

	Is_expired  bool `json:"is_expired"`
	Is_eligible bool `json:"is_eligible"`
	Used        bool `json:"used"`
	Valid       bool
	Present     bool
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
	Cartid       string `json:"cart_id"`
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
		return err, Coupon
	}
	if !(u.Type == "ONLINE" || u.Type == "COD") {
		err.Add("Payment Type  ", "please ADD Payment Type")
		return err, Coupon
	}
	// if u.Aid == "" {
	// 	err.Add("Aid ", "no address")
	// 	return err
	// }
	return url.Values{}, Coupon

}
