package handlers

import (
	"GoProject/middleware"
	"GoProject/models"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

func SeeProfileDetail(c *gin.Context) {
	// Cancel blocking function after 10 sec.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	pathId := c.Param("id")
	ObjectId, errHex := primitive.ObjectIDFromHex(pathId)
	if errHex != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errHex.Error()})
		return
	}

	var visitUser models.User
	filter := bson.D{{"_id", ObjectId}}
	err := userCollection.FindOne(ctx, filter).Decode(&visitUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": errHex.Error()})
		return
	}

	jsonUser, errMarsh := json.Marshal(visitUser)
	if errMarsh != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errMarsh.Error()})
		return
	}
	c.JSON(http.StatusOK, jsonUser)
}

func SeeMyProfile(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	userData := c.Keys["user_data"].(middleware.UserTokenData)
	ObjectId, errHex := primitive.ObjectIDFromHex(userData.Id)
	if errHex != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errHex.Error()})
		return
	}

	var User models.User
	filter := bson.D{{"_id", ObjectId}}
	err := userCollection.FindOne(ctx, filter).Decode(User)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": errHex.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": User})
}
