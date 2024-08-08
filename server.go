package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/manthan1609/authentication-app/database"
	"github.com/manthan1609/authentication-app/routes"
)

func init() {
	fmt.Println("Yhi H")

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln("error in loading environment variables")
	}

	fmt.Println("Yhi H")
	database.ConnectDB()
	fmt.Println("Yhi H")

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
