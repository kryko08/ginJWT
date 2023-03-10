package middleware

import (
	"GoProject/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestExtractTokenString(t *testing.T) {
	utils.LoadEnv("../.env")
	s := JWTService{}
	s.SetUpJWTService("test", *jwt.SigningMethodHS256)
	id := primitive.NewObjectID()
	idS := id.Hex()
	token := s.GenerateJWT(idS)

	req, err := http.NewRequest("GET", "www.example.com", nil)
	if err != nil {
		t.Fatal("Error creating request")
	}
	req.Header.Add("jwt", token)
	// Test function
	tokenString, err := extractTokenString(req)
	if err != nil {
		log.Fatal("Error extracting token")
	}
	log.Printf(tokenString)
}

func TestExtractTokenStringEmpty(t *testing.T) {
	req, err := http.NewRequest("GET", "www.example.com", nil)
	if err != nil {
		t.Fatal("Error creating request")
	}
	// Test function
	_, extractionError := extractTokenString(req)
	log.Print(extractionError)
	if extractionError == nil {
		t.Fatal("No token to be extracted")
	}
}

func TestExtractUserData(t *testing.T) {
	// generate and verify token
	s := JWTService{}
	s.SetUpJWTService("test", *jwt.SigningMethodHS256)

	id := primitive.NewObjectID()
	idS := id.Hex()
	tokenString := s.GenerateJWT(idS)
	token, err := s.VerifyJWT(tokenString)
	if err != nil {
		t.Fatal("error verifying token")
	}
	userData, extractErr := extractUserData(token)
	if extractErr != nil {
		t.Fatal("Error extracting user Data")
	}
	extractIdString := userData.Id
	extractIdObjectId, transformErr := primitive.ObjectIDFromHex(extractIdString)
	if transformErr != nil {
		t.Fatal("Error from hex to objectId")
	}

	if extractIdObjectId != id {
		t.Fatal("Ids do not match")
	}

}

func TestJWTAuthorizationWithToken(t *testing.T) {
	// create service and create JWT token
	s := JWTService{}
	s.SetUpJWTService("test", *jwt.SigningMethodHS256)
	id := primitive.NewObjectID()
	idS := id.Hex()
	token := s.GenerateJWT(idS)

	// setup router with middleware and handler
	r := gin.Default()
	r.GET("test/", JWTAuthorization(s), func(context *gin.Context) {
		// Get header
		resp, ok := context.Get("user_data")
		if !ok {
			t.Fatal("could not get value from context")
		}
		context.JSON(http.StatusOK, gin.H{"response": resp})
	})

	// create request
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	req, err := http.NewRequestWithContext(ctx, "GET", "/test/", nil)
	if err != nil {
		t.Fatal("error creating request")
	}
	// put token into header
	req.Header.Add("jwt", token)
	r.ServeHTTP(w, req)
	// handlers write calls are sent to Body
	responseCode := w.Code
	assert.Equal(t, 200, responseCode, "wrong response code")
}

func TestJWTAuthorizationWithNoToken(t *testing.T) {
	// create service and create JWT token
	s := JWTService{}
	s.SetUpJWTService("test", *jwt.SigningMethodHS256)

	// setup router with middleware and handler
	r := gin.Default()
	r.GET("test/", JWTAuthorization(s), func(context *gin.Context) {
		// Get header
		resp, ok := context.Get("user_data")
		if !ok {
			t.Fatal("could not get value from context")
		}
		context.JSON(http.StatusOK, gin.H{"response": resp})
	})

	// create request
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	req, err := http.NewRequestWithContext(ctx, "GET", "/test/", nil)
	if err != nil {
		t.Fatal("error creating request")
	}
	r.ServeHTTP(w, req)
	// handlers write calls are sent to Body
	responseCode := w.Code
	assert.Equal(t, 401, responseCode, "wrong response code")
}

func TestJWTAuthorizationWithAlteredToken(t *testing.T) {
	// create service and create JWT token
	s := JWTService{}
	s.SetUpJWTService("test", *jwt.SigningMethodHS256)
	id := primitive.NewObjectID()
	idS := id.Hex()
	token := s.GenerateJWT(idS)
	alteredToken := utils.AlterToken(token)

	// setup router with middleware and handler
	r := gin.Default()
	r.GET("test/", JWTAuthorization(s), func(context *gin.Context) {
		// Get header
		resp, ok := context.Get("user_data")
		if !ok {
			t.Fatal("could not get value from context")
		}
		context.JSON(http.StatusOK, gin.H{"response": resp})
	})

	// create request
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	req, err := http.NewRequestWithContext(ctx, "GET", "/test/", nil)
	if err != nil {
		t.Fatal("error creating request")
	}
	// put token into header
	req.Header.Add("jwt", alteredToken)
	r.ServeHTTP(w, req)
	// handlers write calls are sent to Body
	responseCode := w.Code
	assert.Equal(t, 401, responseCode, "wrong response code")
}
