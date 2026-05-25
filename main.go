package main

import (
	"gin-project/config"
	"gin-project/routes"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			env.Getenv("CORS_ALLOWED_ORIGINS"),
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// db connection
	config.ConnectDB()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Gin server working",
		})
	})
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Gin server working",
		})
	})

	routes.Routes(r)

	// r.Use("/api", routes.Routes)

	if err := r.Run(":" + env.Getenv("PORT")); err != nil {
		log.Fatal(err)
	}
}
