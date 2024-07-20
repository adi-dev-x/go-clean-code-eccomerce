package model

import (
	"net/url"

	"github.com/golang-jwt/jwt"
)

type UserRegisterRequest struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Phone     string `json:"phone"`
}
type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type UserClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
type UserOtp struct {
	Email string `json:"email"`
	Otp   string `json:"otp"`
}

func (u *UserRegisterRequest) Valid() url.Values {
	err := url.Values{}

	if len(u.FirstName) < 2 {
		err.Add("first_name", "invalid first name")
	}

	if len(u.Password) < 6 {
		err.Add("password", "password must be greater than 6")
	}

	return err
}
