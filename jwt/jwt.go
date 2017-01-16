package jwt

import (
	"time"

	"github.com/NorbertKa/LambdaCMS/models"
	"github.com/dgrijalva/jwt-go"
)

type TokenInfo struct {
	jwt.StandardClaims
	UserId   int
	Username string
	Role     string
}

func DecodeToken(token string, secret string) (*TokenInfo, error) {
	claims := TokenInfo{}
	_, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	return &claims, nil
}

func EncodeToken(user db.User, secret string) (string, error) {
	tok := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"userId":   user.Id,
		"username": user.Username,
		"role":     user.Role,
		"nbf":      time.Now().Unix(),
	})
	tokenString, err := tok.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
