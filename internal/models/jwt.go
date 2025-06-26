package models

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"google.golang.org/api/idtoken"
)

const GoogleId = "652300787712-178nsm16d8e7o6ia6a763c5unjvhudss.apps.googleusercontent.com"

func GetJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "secret"
	}
	return secret
}

var (
	JWT_SECRET          = GetJWTSecret()
	ErrInvalidToken     = errors.New("invalid token")
	ErrExpiredToken     = errors.New("token expired")
	ErrMissingDataField = errors.New("missing data field")
	ErrJWTValidation    = jwt.NewValidationError("invalid signing method", jwt.ValidationErrorSignatureInvalid)
)

func ValidateGoogleToken(ctx context.Context, token string) (CreateUserRequest, error) {
	var ok bool
	user := CreateUserRequest{
		Role:     "student",
		Verified: true,
	}

	payload, err := idtoken.Validate(ctx, token, GoogleId)
	if err != nil {
		return user, fmt.Errorf("failed to validate token: %w", err)
	}

	user.Email, ok = payload.Claims["email"].(string)
	if !ok {
		return user, fmt.Errorf("invalid token: missing email claim")
	}

	user.Name, ok = payload.Claims["given_name"].(string)
	if !ok {
		return user, fmt.Errorf("invalid token: missing given_name claim")
	}

	user.Surname, ok = payload.Claims["family_name"].(string)
	if !ok {
		return user, fmt.Errorf("invalid token: missing family_name claim")
	}

	return user, nil
}

type Claims struct {
	jwt.StandardClaims
	Email string `json:"email"`
	Name  string `json:"full_name"`
	Role  string `json:"role"`
	Admin bool   `json:"admin"`
}

// GenerateToken genera un token JWT para el usuario.
func GenerateToken(id int, email string, name string, role string) (string, error) {
	claims := Claims{
		StandardClaims: jwt.StandardClaims{
			Subject:   strconv.Itoa(id),
			Issuer:    "user-api",
			ExpiresAt: time.Now().Add(time.Hour * 24 * 7).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		Email: email,
		Name:  name,
		Role:  role,
	}
	if role == "admin" {
		claims.Admin = true
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(GetJWTSecret()))
}

func ParseToken(tokenStr string) (*Claims, error) {
	secret := []byte(JWT_SECRET)
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
