package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/manthan1609/authentication-app/controllers"
)

func AuthRouter(app *gin.Engine) {
	app.POST("/auth/signup", controller.SignUp)
	app.POST("/auth/signin", controller.SignIn)
}
