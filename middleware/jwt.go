package middleware

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"time"
)

type TokenGeneratorVerifier interface {
	generateJWT(string2 string) string
	verifyJWT(string2 string) (*jwt.Token, string)
}

// JWTService Implements TokenGeneratorVerifier
type JWTService struct {
	publicKey  string
	privateKey string
	signMethod jwt.SigningMethod
	issuer     string
}

type customClaims struct {
	Id string `json:"id"`
	jwt.RegisteredClaims
}

func SetUpJWTService(signMethod jwt.SigningMethod, issuer string, keys ...string) (JWTService, error) {
	switch signMethod.(type) {
	case *jwt.SigningMethodECDSA, *jwt.SigningMethodRSA:
		return JWTService{
			privateKey: keys[0],
			publicKey:  keys[1],
			signMethod: signMethod,
			issuer:     issuer,
		}, nil

	case *jwt.SigningMethodHMAC:
		return JWTService{
			privateKey: keys[0],
			publicKey:  "",
			signMethod: signMethod,
			issuer:     issuer,
		}, nil

	default:
		return JWTService{}, errors.New("not supported signing method as parameter")
	}
}

func (s *JWTService) isAsymmetricSign() bool {
	switch s.signMethod.(type) {
	case *jwt.SigningMethodECDSA, *jwt.SigningMethodRSA:
		return true
	}
	return false
}

func (s *JWTService) GenerateJWT(userId string) string {
	claims := customClaims{
		userId,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "kryko08",
		},
	}
	token := jwt.NewWithClaims(s.signMethod, claims)
	// use private key both for symmetric and asymmetric token generation
	t, err := token.SignedString([]byte(s.privateKey))
	if err != nil {
		log.Fatal("Error signing token", err)
	}
	return t
}

func (s *JWTService) VerifyJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// validate alg
		if _, ok := token.Method.(jwt.SigningMethod); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		if s.isAsymmetricSign() {
			return []byte(s.publicKey), nil
		} else {
			return []byte(s.privateKey), nil
		}
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}
