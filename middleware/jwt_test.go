package middleware

import (
	"GoProject/utils"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestGenerateJWT(t *testing.T) {
	service := JWTService{}
	service.SetUpJWTService("test", *jwt.SigningMethodHS256)
	id := primitive.NewObjectID()
	idS := id.Hex()
	token := service.GenerateJWT(idS)
	_, err := service.VerifyJWT(token)
	if err != nil {
		t.Fatal("Error parsing the token", err)
	}
}

func TestAlterToken(t *testing.T) {
	service := JWTService{}
	service.SetUpJWTService("test", *jwt.SigningMethodHS256)
	id := primitive.NewObjectID()
	idS := id.Hex()
	token := service.GenerateJWT(idS)
	alteredToken := utils.AlterToken(token)
	_, err := service.VerifyJWT(alteredToken)
	if err == nil {
		t.Fatal("Expected error, Token altered")
	}
}
