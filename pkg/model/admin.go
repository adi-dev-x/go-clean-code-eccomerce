package model

import (
	"net/url"

	"github.com/golang-jwt/jwt"
)

type AdminRegisterRequest struct {
	Name string `json:"name"`

	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	GST      string `json:"gst"`
}
type AdminLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type AdminClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
type AdminOtp struct {
	Email string `json:"email"`
	Otp   string `json:"otp"`
}

func (u *AdminRegisterRequest) Valid() url.Values {
	err := url.Values{}

	if len(u.Password) < 6 {
		err.Add("password", "password must be greater than 6")
	}

	return err
}
