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
			"http://localhost:3000",
			"http://localhost:3001",
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

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
