package main

import (
	"GoProject/db"
	"GoProject/handlers"
	"GoProject/middleware"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"net/http"
)

func RouterSetUp() *gin.Engine {
	// jwt setup
	var jwtS = middleware.JWTService{}
	jwtS.SetUpJWTService("ginJWT", *jwt.SigningMethodHS256)
	var loginS = handlers.LoginService{
		JWT: jwtS,
	}

	r := gin.Default()

	db.ConnectDB()

	auth := r.Group("/auth")
	{
		auth.POST("/login", func(c *gin.Context) {
			token, err := loginS.Login(c)
			if err.Err != nil {
				c.JSON(err.StatusCode, gin.H{"error": err.Error()})
				return
			}
			// add token to header
			c.Header("jwt", token)
			c.JSON(http.StatusOK, gin.H{"success": "jwt token added to response header"})
			return
		})
		auth.POST("/register", handlers.RegisterUser)
	}

	users := r.Group("/users")
	{
		users.GET("/:id", handlers.SeeProfileDetail)
	}

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"response": "pong"})

	})

	return r
}

func main() {
	router := RouterSetUp()
	// run server
	err := router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	if err != nil {
		log.Fatal("Cannot run router", err)
	}
}
