package main

import (
	"go-blog/config"
	"go-blog/database"
	"go-blog/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()
	database.ConnectDB()

	r := gin.Default()
	routes.RegisterRoutes(r)

	// port := os.Getenv("PORT")
	r.Run("0.0.0.0:8080")
}
