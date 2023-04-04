package auth

import (
	"errors"
	"os"

	"github.com/dgrijalva/jwt-go"
)

type AuthService interface {
	GenerateToken(userID int) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}

type authService struct{}

var SECREAT_KEY = []byte(os.Getenv("API_SECRET"))

func NewAuthService() *authService {
	return &authService{}
}

func (s *authService) ValidateToken(encodedToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("Invalid Token")
		}
		return []byte(SECREAT_KEY), nil
	})
	if err != nil {
		return token, err
	}
	return token, nil
}

func (s *authService) GenerateToken(userID int) (string, error) {
	claim := jwt.MapClaims{}
	claim["user_id"] = userID

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	signedToken, err := token.SignedString([]byte(SECREAT_KEY))
	if err != nil {
		return signedToken, err
	}
	return signedToken, nil
}
