package middleware

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type MyClaims struct {
	Id string `json:"id"`
	jwt.RegisteredClaims
}

func getMySecret() string {

}

func generateJWT(userId string) (string, error) {
	claims := MyClaims{
		userId,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// add claims

}
