package admin

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

type AdminJWT interface {
	GenerateAdminToken(username string) (string, error)
	AdminAuthMiddleware() echo.MiddlewareFunc
}

type Adminjwt struct {
	Config config.Config
}

func (s Adminjwt) GenerateAdminToken(username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := model.AdminClaims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(s.Config.VnJWTKey))
}

// verify Admin Token
func AdminAuthentication(tokenString string, jwtKey string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &model.AdminClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*model.AdminClaims); ok && token.Valid {
		return claims.Username, nil
	}
	return "", errors.New("invalid token")
}

// Admin Auth middleware
func (s Adminjwt) AdminAuthMiddleware() echo.MiddlewareFunc {
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

			username, err := AdminAuthentication(tokenString, s.Config.VnJWTKey)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token"})
			}
			c.Set("username", username)
			return next(c)
		}
	}
}
