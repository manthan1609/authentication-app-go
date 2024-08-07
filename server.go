package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/manthan1609/authentication-app/database"
	"github.com/manthan1609/authentication-app/routes"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln("error in loading environment variables")
	}

	database.ConnectDB()
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	app := gin.New()
	app.Use(gin.Logger())

	routes.AuthRouter(app)
	routes.UserRouter(app)

	app.Run(":" + port)
}
