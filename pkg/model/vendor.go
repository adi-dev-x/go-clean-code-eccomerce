package model

import (
	"net/url"

	"github.com/golang-jwt/jwt"
)

type VendorRegisterRequest struct {
	Name string `json:"name"`

	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	GST      string `json:"gst"`
}
type VendorLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type VendorClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
type VendorOtp struct {
	Email string `json:"email"`
	Otp   string `json:"otp"`
}

func (u *VendorRegisterRequest) Valid() url.Values {
	err := url.Values{}

	if len(u.Password) < 6 {
		err.Add("password", "password must be greater than 6")
	}

	return err
}
