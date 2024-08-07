package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/manthan1609/authentication-app/database"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.OpenCollection("user")
var validate = validator.New()

func GetUsers(c *gin.Context) {

}

func GetUser(c *gin.Context) {

}
