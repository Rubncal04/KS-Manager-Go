package middleware

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func Authentication(c echo.Context) *JWTClaims {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	jwtClaims := JWTClaims{
		UserID:   uint(claims["user_id"].(float64)),
		Name:     claims["name"].(string),
		UserName: claims["user_name"].(string),
		ChurchId: int(claims["church_id"].(float64)),
		Role:     claims["role"].(string),
	}

	return &jwtClaims
}
