package controllers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/manthan1609/authentication-app/database"
	"github.com/manthan1609/authentication-app/helpers"
	"github.com/manthan1609/authentication-app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var validate = validator.New()

func GetUsers(c *gin.Context) {

	var userCollection *mongo.Collection = database.OpenCollection("user")

	if err := helpers.CheckUserType(c, "ADMIN"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))

	if err != nil || recordPerPage < 1 {
		recordPerPage = 10
	}

	page, err := strconv.Atoi(c.Query("page"))

	if err != nil || page < 1 {
		page = 1
	}

	startIndex, err := strconv.Atoi(c.Query("startIndex"))

	if err != nil {
		startIndex = page * recordPerPage
	}

	matchStage := bson.D{{Key: "$match", Value: bson.D{{}}}}

	groupStage := bson.D{
		{
			Key: "$group",
			Value: bson.D{
				{Key: "_id", Value: bson.D{{Key: "_id", Value: "null"}}},
				{Key: "total_count", Value: bson.D{{Key: "$sum", Value: 1}}},
				{Key: "data", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}},
			},
		},
	}

	projectStage := bson.D{{
		Key: "$project",
		Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "total_count", Value: 1},
			{Key: "user_items", Value: bson.D{
				{Key: "$slice", Value: []interface{}{"$data", startIndex}},
			}},
		},
	}}

	result, err := userCollection.Aggregate(ctx, mongo.Pipeline{
		matchStage,
		groupStage,
		projectStage,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "error occured while fetching users",
			"err":   err,
		})
		return
	}

	var allUsers []bson.M

	if err := result.All(ctx, &allUsers); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "error occured while fetching users",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": allUsers[0],
	})
}

func GetUser() gin.HandlerFunc {
	var userCollection *mongo.Collection = database.OpenCollection("user")

	return func(c *gin.Context) {
		userId := c.Param("id")

		if err := helpers.MatchUserTypeToUid(c, userId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User

		err := userCollection.FindOne(ctx, bson.M{
			"user_id": userId,
		}).Decode(&user)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}
