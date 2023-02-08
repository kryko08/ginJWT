package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"net/http"
	"strings"
)

// UserTokenData stores user data after token has been verified
type UserTokenData struct {
	Id string
}

func extractTokenString(h *http.Request) (string, error) {
	header := h.Header
	tokenString := header.Get("jwt") // get jwt token string from http header
	// check if header does not exist
	if len(strings.TrimSpace(tokenString)) == 0 {
		return " ", errors.New("token not found")
	}
	return tokenString, nil
}

func extractUserData(token *jwt.Token) (*UserTokenData, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok {
		userId, userIdOk := claims["id"].(string)
		if userIdOk == false {
			return nil, errors.New("could not get User id from token Claims")
		} else {
			utd := &UserTokenData{
				Id: userId,
			}
			return utd, nil
		}
	}
	return nil, errors.New("could not map claims")
}

func ValidateRequest(ctx *gin.Context, s *JWTService) (*UserTokenData, error) {
	r := ctx.Request
	tokenString, errExtract := extractTokenString(r)
	if errExtract != nil {
		return nil, errors.New("could not extract token string")
	}
	token, errVerif := s.VerifyJWT(tokenString)
	if errVerif != nil {
		return nil, errors.New("token not verified")
	}
	userData, errData := extractUserData(token)
	if errData != nil {
		return nil, errors.New("could not extract data")
	}
	return userData, nil
}

func JWTAuthorization(s JWTService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userData, err := ValidateRequest(ctx, &s)
		if err != nil {
			log.Print("TADY JE ERROR", err)
			ctx.JSON(http.StatusUnauthorized, err)
			ctx.Abort()
			return
		}
		ctx.Set("user_data", userData)
		ctx.Next()
	}
}
