package authorization

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claim struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"correo"`
	jwt.RegisteredClaims
}

// recibir el login y generar un toekn
func GenerateToken(userID int64, email string) (string, error) {
	claims := &Claim{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Minute)),
			Issuer:    "api-module",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// generar el token a partir de los claims y preparar el token para firmarlo con la clave privada
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	// firmamos el token con la clave privada y retronar
	return token.SignedString(signKey)
}

// validar el token
func ValidateToken(tokenString string) (*Claim, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claim{}, func(token *jwt.Token) (interface{}, error) {
		// validar el método de la firma sea RSA
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return verifyKey, nil
	})
	if err != nil {
		return nil, err
	}

	// verificar la validez de la firma y la fecha de expiración
	if !token.Valid {
		return nil, errors.New("the token is invalid or has expired")
	}

	claims, ok := token.Claims.(*Claim)
	if !ok {
		return nil, errors.New("the token claims could not be processed")
	}

	return claims, nil
}
