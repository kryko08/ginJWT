package handlers

import (
	"GoProject/db"
	"GoProject/middleware"
	"GoProject/models"
	"GoProject/utils"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

var userCollection = db.GetCollection("users", db.Database)

type LoginService struct {
	JWT middleware.JWTService
}

func RegisterUser(c *gin.Context) {
	// Cancel blocking function after 10 sec.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var newUser models.PostUser
	bindErr := c.BindJSON(&newUser)
	if bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": bindErr.Error()})
		return
	}

	// check if passwords match
	if newUser.Password != newUser.PasswordConfirm {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Passwords do not match"})
		return
	}

	// hash password
	password := newUser.Password
	passwordHash, errHash := utils.HashPassword(password)
	if errHash != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errHash.Error()})
		return
	}

	newId := primitive.NewObjectID()
	userToSave := models.User{
		ID:       newId,
		Username: newUser.Username,
		Email:    newUser.Email,
		Password: passwordHash,
	}

	_, errInsert := userCollection.InsertOne(ctx, userToSave)
	if errInsert != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errInsert.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"StatusOk": "User has been successfully added to database"})
}

func (l *LoginService) Login(c *gin.Context) (string, utils.RequestError) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var loginDetails models.LoginForm
	bindErr := c.ShouldBindJSON(&loginDetails)
	if bindErr != nil {
		return "", utils.RequestError{
			StatusCode: 400,
			Err:        errors.New("could not bind JSON"),
		}
	}

	var user models.User
	filter := bson.D{{"username", loginDetails.Username}}
	findErr := userCollection.FindOne(ctx, filter).Decode(&user)
	if findErr != nil {
		if findErr == mongo.ErrNoDocuments {
			return "", utils.RequestError{
				StatusCode: 200,
				Err:        errors.New("no user with this username"),
			}
		}
		return "", utils.RequestError{
			StatusCode: 500,
			Err:        errors.New("error fetching user data from database"),
		}
	}

	// password check
	check := utils.CheckPassword(loginDetails.Password, user.Password)
	if check != true {
		return "", utils.RequestError{
			StatusCode: 401,
			Err:        errors.New("incorrect password"),
		}
	}

	// create token
	userId := user.ID.Hex()
	token := l.JWT.GenerateJWT(userId)
	return token, utils.RequestError{
		StatusCode: 200,
		Err:        nil,
	}
}
