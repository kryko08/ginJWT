package middleware

import (
	"GoProject/utils"
	"github.com/golang-jwt/jwt/v4"
	"testing"
)

func TestGenerateJWT(t *testing.T) {
	service := JWTService{}
	service.setUpJWTService("test", *jwt.SigningMethodHS256)
	token := service.generateJWT("1")
	_, err := service.verifyJWT(token)
	if err != nil {
		t.Fatal("Error parsing the token", err)
	}
}

func TestAlterToken(t *testing.T) {
	service := JWTService{}
	service.setUpJWTService("test", *jwt.SigningMethodHS256)
	token := service.generateJWT("1")
	alteredToken := utils.AlterToken(token)
	_, err := service.verifyJWT(alteredToken)
	if err == nil {
		t.Fatal("Expected error, Token altered")
	}
}
