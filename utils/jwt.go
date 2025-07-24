package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func GenerateToken(userId int) (string, error) {
	godotenv.Load()

	expirationTime := time.Now().Add(12 * time.Hour)
	claims := jwt.MapClaims{
		"userId": userId,
		"iat":    time.Now().Unix(),
		"exp":    expirationTime.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secretKey := os.Getenv("SECRET_KEY")
	return token.SignedString([]byte(secretKey))
}