package middleware

import (
	"time"

	"github.com/Rubncal04/ksmanager/config"
	"github.com/Rubncal04/ksmanager/models"
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Name     string `json:"name"`
	UserName string `json:"user_name"`
	ChurchId int    `json:"church_id"`
	Role     int    `json:"role"`
	jwt.RegisteredClaims
}

func GenerateJWT(user *models.User) (string, error) {
	variables := config.GetVariables()
	expirationTime := time.Now().Add(72 * time.Hour)
	var jwtKey = []byte(variables.JWT_KEY)

	claims := &JWTClaims{
		UserID:   user.ID,
		UserName: user.UserName,
		Name:     user.Name,
		ChurchId: user.ChurchId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
