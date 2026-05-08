package main

import (
	"gin-project/config"
	"gin-project/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

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
