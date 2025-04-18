package utils

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// GenerateToken genera un token JWT para el usuario.
func GenerateToken(userId int) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "secret" // valor por defecto
	}

	claims := jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
