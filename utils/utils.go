package utils

import (
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"strings"
)

type EnvStr struct {
	MongoUri string
	Key      string
}

var EnvVars = GetEnvVars("./.env")

func GetEnvVars(relativeEnvPath string) *EnvStr {
	err := godotenv.Load(relativeEnvPath)
	if err != nil {
		log.Fatal("Error while loading .env file", err)
	}
	return &EnvStr{
		MongoUri: os.Getenv("MONGO_URI"),
		Key:      os.Getenv("SECRET_KEY"),
	}
}

func AlterToken(s string) string {
	sList := strings.Split(s, ".")
	reversedClaims := reverseString(sList[1])
	stringArray := []string{sList[0], reversedClaims, sList[2]}
	alteredToken := strings.Join(stringArray, ".")
	return alteredToken
}

func reverseString(str string) string {
	byteStr := []rune(str)
	for i, j := 0, len(byteStr)-1; i < j; i, j = i+1, j-1 {
		byteStr[i], byteStr[j] = byteStr[j], byteStr[i]
	}
	return string(byteStr)
}

func HashPassword(p string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(p), 14)
	return string(bytes), err
}

func CheckPassword(p string, h string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(h), []byte(p))
	return err == nil
}
