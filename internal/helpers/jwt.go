package helpers

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const TOKEN_EXP = time.Hour * 3

type Claims struct {
	jwt.RegisteredClaims
	UserLogin string
}

func getJwtSecret() (string, error) {
	secret, exists := os.LookupEnv("AUTH_SECRET")
	if !exists {
		return "", fmt.Errorf("secret not found")
	}
	return secret, nil
}

func CreateJWTToken(login string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			// когда создан токен
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TOKEN_EXP)),
		},
		// собственное утверждение
		UserLogin: login,
	})
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

func GetUserLogin(tokenString string) string {
	// создаём экземпляр структуры с утверждениями
	claims := &Claims{}
	// парсим из строки токена tokenString в структуру claims
	jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		secretKey, err := getJwtSecret()
		if err != nil {
			return "", err
		}
		return []byte(secretKey), nil
	})

	// возвращаем ID пользователя в читаемом виде
	return claims.UserLogin
}
