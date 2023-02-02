package utils

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"
)

func GetOsEnv(osEnv string) string {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error while loading .env file", err)
	}
	secret := os.Getenv(osEnv)
	if secret == "" {
		log.Fatal("error loading a secret", secret)
	}
	return secret
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
