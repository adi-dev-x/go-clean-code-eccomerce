package user

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"myproject/pkg/config"
	"myproject/pkg/model"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

// AdminJWT interface defines methods for JWT operations
type AdminJWT interface {
	GenerateAdminToken(username string) (string, error)
	AdminAuthMiddleware() echo.MiddlewareFunc
}

// Adminjwt struct holds configuration
type Adminjwt struct {
	Config config.Config
}

// GenerateAdminToken generates a JWT token for admin users
func (s Adminjwt) GenerateAdminToken(username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := model.UserClaims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(s.Config.AdJWTKey))
}

// AdminAuthentication verifies the admin token
func AdminAuthentication(tokenString, jwtKey string) (string, error) {
	fmt.Println("Inside AdminAuthentication")
	token, err := jwt.ParseWithClaims(tokenString, &model.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*model.UserClaims); ok && token.Valid {
		return claims.Username, nil
	}
	return "", errors.New("invalid token")
}

// AdminAuthMiddleware provides middleware for admin authentication
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

			username, err := AdminAuthentication(tokenString, s.Config.AdJWTKey)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token"})
			}
			c.Set("username", username)
			return next(c)
		}
	}
}
