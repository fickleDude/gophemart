package helpers

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v4"
)

// change to os env
func getJwtSecret() (string, error) {
	secret, exists := os.LookupEnv("AUTH_SECRET")
	if exists {
		return "", fmt.Errorf("secret not found")
	}
	return secret, nil
}

func CreateJWTToken(login string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"login": login})
	secretKey, err := getJwtSecret()
	if err != nil {
		return "", err
	}
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidateJWTToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		secretKey, err := getJwtSecret()
		if err != nil {
			return "", err
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return nil
	}
	return jwt.ErrSignatureInvalid
}
