package models

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	ErrInvalidToken     = errors.New("invalid token")
	ErrExpiredToken     = errors.New("token expired")
	ErrMissingDataField = errors.New("missing data field")
	ErrJWTValidation    = jwt.NewValidationError("invalid signing method", jwt.ValidationErrorSignatureInvalid)
)

type Claims struct {
	jwt.StandardClaims
	Email string `json:"email"`
	Name  string `json:"full_name"`
	Role  string `json:"role"`
}

func GetJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "secret" // valor por defecto
	}
	return secret
}

// GenerateToken genera un token JWT para el usuario.
func GenerateToken(id int, email string, name string, role string) (string, error) {
	claims := Claims{
		StandardClaims: jwt.StandardClaims{
			Subject:   strconv.Itoa(id),
			Issuer:    "user-api",
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		Email: email,
		Name:  name,
		Role:  role,
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(GetJWTSecret()))
}

func ParseToken(tokenStr string) (*Claims, error) {
	secret := []byte(GetJWTSecret())
	var claims Claims
	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrJWTValidation
		}
		return secret, nil
	})

	if err != nil || !token.Valid {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, ErrInvalidToken
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, ErrExpiredToken
			} else {
				return nil, ErrJWTValidation
			}
		}
		return nil, ErrInvalidToken
	}

	return &claims, nil
}
