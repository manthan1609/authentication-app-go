package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/manthan1609/authentication-app/database"
	"github.com/manthan1609/authentication-app/helpers"
	"github.com/manthan1609/authentication-app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func SignUp() gin.HandlerFunc {

	return func(c *gin.Context) {
		var userCollection *mongo.Collection = database.OpenCollection("user")

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid data format",
			})
			return
		}

		validationError := validate.Struct(user)
		if validationError != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": validationError.Error(),
			})
			return
		}

		count, err := userCollection.CountDocuments(ctx, bson.M{
			"email": user.Email,
		})

		if err != nil {
			// log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "error while checking for user",
			})
			return
		}

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "user already exist",
			})
			return
		}

		count, err = userCollection.CountDocuments(ctx, bson.M{
			"phone": user.Phone,
		})

		if err != nil {
			// log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "error while checking for user",
			})
			return
		}

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "user already exist",
			})
			return
		}

		user.CreatedAt, err = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		if err != nil {
			// log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "error while registering user",
			})
			return
		}

		user.UpdatedAt, err = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		if err != nil {
			// log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "error while registering user",
			})
			return
		}
		user.ID = primitive.NewObjectID()
		user.UserId = user.ID.Hex()

		hashedPassword, err := helpers.HashPassword(*user.Password)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "error while registering user",
			})
			return
		}

		user.Password = &hashedPassword

		accessToken, refreshToken, err := helpers.GenerateAllTokens(*user.Email, *user.FirstName, *user.LastName, *user.UserType, user.UserId)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "error while registering user",
			})
			return
		}

		user.Token = &accessToken
		user.RefreshToken = &refreshToken

		resultInsertionNumber, err := userCollection.InsertOne(ctx, user)

		if err != nil {
			// log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "database error while registering user",
			})
			return
		}

		c.SetCookie("access_token", accessToken, 60*60*24, "", "", true, true)
		c.SetCookie("refresh_token", refreshToken, 60*60*24*365, "", "", true, true)

		c.JSON(http.StatusOK, gin.H{
			"message":          "user registered successfully",
			"insertion_number": resultInsertionNumber,
		})
	}

}

func SignIn(c *gin.Context) {
	var userCollection *mongo.Collection = database.OpenCollection("user")

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var user models.User
	var foundUser models.User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid data format",
		})
		return
	}

	validationError := validate.Struct(user)
	if validationError != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": validationError.Error(),
		})
		return
	}

	err := userCollection.FindOne(ctx, bson.M{
		"email": user.Email,
	}).Decode(&foundUser)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "email or password is incorrect",
		})
		return
	}

	if !helpers.VerifyPassword(*foundUser.Password, *user.Password) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "email or password is incorrect",
		})
		return
	}

	if foundUser.Email == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "user not found",
		})
	}

	foundUser.UpdatedAt, err = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	if err != nil {
		// log.Panic(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "error while login the user",
		})
		return
	}

	accessToken, refreshToken, err := helpers.GenerateAllTokens(*foundUser.Email, *foundUser.FirstName, *foundUser.LastName, *foundUser.UserType, foundUser.UserId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "error while login user",
		})
		return
	}

	foundUser.Token = &accessToken
	foundUser.RefreshToken = &refreshToken

	updateResult, err := userCollection.UpdateOne(ctx, bson.M{
		"user_id": foundUser.UserId,
	}, bson.M{
		"$set": bson.M{
			"updated_at":    foundUser.UpdatedAt,
			"token":         foundUser.Token,
			"refresh_token": foundUser.RefreshToken,
		},
	})

	if err != nil {
		// log.Panic(err)
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "database error while log in the user",
		})
		return
	}

	c.SetCookie("access_token", accessToken, 60*60*24, "", "", true, true)
	c.SetCookie("refresh_token", refreshToken, 60*60*24*365, "", "", true, true)

	c.JSON(http.StatusOK, gin.H{
		"message":      "user logged in successfully",
		"updateResult": updateResult,
	})
}
