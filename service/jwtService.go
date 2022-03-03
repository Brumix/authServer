package service

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"os"
	"time"
)

type JWTService interface {
	GenerateToken(email string, userName string) string
	ValidateToken(token string) (*jwt.Token, error)
}
type authCustomClaims struct {
	Name     string `json:"name"`
	UserName string `json:"user_name"`
	jwt.StandardClaims
}

type JwtServices struct{}

func NewJWTService() JwtServices {
	return JwtServices{}
}

func getSecretKey() string {
	secret := os.Getenv("JWTKEY")
	if secret == "" {
		secret = "7a06cf5ff9d3583a483607ae64c9275712662884b7be21ab31b4776e29e7e5ffe7ad889eff44972b38732be32e212c9d1c41892c03b8adb70061808081e60431"
	}
	return secret
}

func (service *JwtServices) GenerateToken(email string, userName string) string {
	claims := &authCustomClaims{
		email,
		userName,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 48).Unix(),
			Issuer:    "Bruno Pereira",
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//encoded string
	t, err := token.SignedString([]byte(getSecretKey()))
	if err != nil {
		panic(err)
	}
	return t
}

func (service *JwtServices) ValidateToken(encodedToken string) (*jwt.Token, error) {
	return jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		if _, isvalid := token.Method.(*jwt.SigningMethodHMAC); !isvalid {
			return nil, fmt.Errorf("invalid token %v", token.Header["alg"])

		}
		return []byte(getSecretKey()), nil
	})

}
