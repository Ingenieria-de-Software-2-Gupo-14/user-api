package auth

import (
	"errors"
	"fmt"
	"ing-soft-2-tp1/internal/models"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	ErrInvalidToken     = errors.New("invalid token")
	ErrExpiredToken     = errors.New("token expired")
	ErrMissingDataField = errors.New("missing data field")
	ErrJWTValidation    = jwt.NewValidationError("invalid signing method", jwt.ValidationErrorSignatureInvalid)
)

type JWTInfo struct {
	UserId int    `json:"user_id"`
	Email  string `json:"email"`
	Admin  bool   `json:"admin"`
	Exp    int64  `json:"exp"`
	Iat    int64  `json:"iat"`
}

func NewJWTInfoFromClaims(claims jwt.MapClaims) (JWTInfo, error) {
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return JWTInfo{}, fmt.Errorf("user_id claim: %v, %w", claims["user_id"], ErrMissingDataField)
	}

	admin, ok := claims["admin"].(bool)
	if !ok {
		return JWTInfo{}, fmt.Errorf("admin claim: %v, %w", claims["admin"], ErrMissingDataField)
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return JWTInfo{}, fmt.Errorf("exp claim: %v, %w", claims["exp"], ErrMissingDataField)
	}

	iat, ok := claims["iat"].(float64)
	if !ok {
		return JWTInfo{}, fmt.Errorf("iat claim: %v, %w", claims["iat"], ErrMissingDataField)
	}

	return JWTInfo{
		UserId: int(userID),
		Admin:  admin,
		Exp:    int64(exp),
		Iat:    int64(iat),
	}, nil
}

func GetJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "secret" // valor por defecto
	}
	return secret
}

// GenerateToken genera un token JWT para el usuario.
func GenerateToken(user models.User) (string, error) {
	secret := GetJWTSecret()

	claims := jwt.MapClaims{
		"user_id": user.Id,
		"email":   user.Email,
		"admin":   user.Admin,
		"iat":     time.Now().Unix(),
		"exp":     time.Now().Add(time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ParseToken(tokenStr string) (JWTInfo, error) {
	secret := []byte(GetJWTSecret())
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrJWTValidation
		}
		return secret, nil
	})

	if err != nil || !token.Valid {
		return JWTInfo{}, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return JWTInfo{}, ErrInvalidToken
	}

	jwtInfo, err := NewJWTInfoFromClaims(claims)
	if err != nil {
		return JWTInfo{}, err
	}

	if jwtInfo.Exp < time.Now().Unix() {
		return JWTInfo{}, ErrExpiredToken
	}

	return jwtInfo, nil
}
