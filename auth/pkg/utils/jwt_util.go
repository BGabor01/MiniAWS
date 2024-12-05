package utils

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte(os.Getenv("JWT_KEY"))

type Claims struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	ID    uint   `json:"userId"`
	jwt.StandardClaims
}

func GenerateJWT(email string, name string, id uint) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Email: email,
		Name: name,
		ID: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}