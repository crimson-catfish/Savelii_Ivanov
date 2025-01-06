package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const expirationTime = time.Hour

var jwtSecret = []byte("very-secret-key-nobody-would-ever-think-of")

func CreateJWTToken(username string) (string, error) {
	claims := jwt.MapClaims{
		"name": username,
		"exp":  time.Now().UTC().Add(expirationTime).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func KeyFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}
	return jwtSecret, nil
}
