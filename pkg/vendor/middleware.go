package vendor

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"myproject/pkg/config"
	"myproject/pkg/model"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type VendorJWT interface {
	GenerateVendorToken(username string) (string, error)
	VendorAuthMiddleware() echo.MiddlewareFunc
}

type Vendorjwt struct {
	Config config.Config
}

func (s Vendorjwt) GenerateVendorToken(username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := model.VendorClaims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(s.Config.VnJWTKey))
}

// verify Vendor Token
func VendorAuthentication(tokenString string, jwtKey string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &model.VendorClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*model.VendorClaims); ok && token.Valid {
		return claims.Username, nil
	}
	return "", errors.New("invalid token")
}

// Vendor Auth middleware
func (s Vendorjwt) VendorAuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing the authorization header"})
			}

			tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer"))
			if tokenString == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid authorization header format"})
			}

			username, err := VendorAuthentication(tokenString, s.Config.VnJWTKey)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token"})
			}
			c.Set("username", username)
			return next(c)
		}
	}
}
