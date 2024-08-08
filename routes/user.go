package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/manthan1609/authentication-app/controllers"
	"github.com/manthan1609/authentication-app/middlewares"
)

func UserRouter(app *gin.Engine) {
	app.Use(middlewares.Authenticate())
	app.GET("/user/:id", controllers.GetUser())
	app.GET("/users", controllers.GetUsers)
}
