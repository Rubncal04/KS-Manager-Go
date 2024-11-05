package middleware

import (
	"fmt"
	"time"

	"github.com/Rubncal04/ksmanager/config"
	"github.com/Rubncal04/ksmanager/db"
	"github.com/Rubncal04/ksmanager/models"
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Name     string `json:"name"`
	UserName string `json:"user_name"`
	ChurchId int    `json:"church_id"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateJWT(user *models.User, repo db.PostgresRepo) (string, error) {
	variables := config.GetVariables()
	expirationTime := time.Now().Add(168 * time.Hour)
	var jwtKey = []byte(variables.JWT_KEY)
	roleId := fmt.Sprintf("%d", user.RoleId)
	role, err := repo.FindOneRoleById(roleId)

	if err != nil {
		return "", err
	}

	claims := &JWTClaims{
		UserID:   user.ID,
		UserName: user.UserName,
		Name:     user.Name,
		ChurchId: user.ChurchId,
		Role:     role.Name,
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
