package middleware

import (
	"GoProject/utils"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"time"
)

type TokenGeneratorVerifier interface {
	generateJWT(string) string
	verifyJWT(string) (*jwt.Token, string)
}

// JWTService Implements TokenGeneratorVerifier
type JWTService struct {
	secret string
	issuer string
	method jwt.SigningMethodHMAC
}

type customClaims struct {
	Id string `json:"id"`
	jwt.RegisteredClaims
}

func (s *JWTService) setUpJWTService(issuer string, hmac jwt.SigningMethodHMAC) {
	s.method = hmac
	// get secret key
	s.secret = utils.GetOsEnv("SECRET_KEY")
	// get issuer
	s.issuer = issuer
}

func (s *JWTService) generateJWT(userId string) string {
	claims := customClaims{
		userId,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    s.issuer,
		},
	}
	// Use symmetric key encryption for signing and validating JWT || Use asymmetric private key for signing and
	// public key for validation
	token := jwt.NewWithClaims(&s.method, claims)
	t, err := token.SignedString([]byte(s.secret))
	if err != nil {
		log.Fatal("Error signing token", err)
	}
	return t
}

func (s *JWTService) verifyJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// validate alg
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secret), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}
