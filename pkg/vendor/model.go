package vendor

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
type AdminClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
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
