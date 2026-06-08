package main

import (
	"gin-project/config"
	"gin-project/queue"
	"gin-project/routes"
	"gin-project/workers"
	"log"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func corsOrigins() []string {
	raw := envOrDefault("CORS_ALLOWED_ORIGINS", "http://localhost:8081")
	parts := strings.Split(raw, ",")
	origins := make([]string, 0, len(parts))
	for _, o := range parts {
		if trimmed := strings.TrimSpace(o); trimmed != "" {
			origins = append(origins, trimmed)
		}
	}
	if len(origins) == 0 {
		return []string{"http://localhost:8081"}
	}
	return origins
}

func main() {
	_ = godotenv.Load()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     corsOrigins(),
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// db connection
	config.ConnectDB()

	// redis connection
	config.InitRedis()

	// rabbitmq connection
	config.InitRabbitMQ()
	defer config.CloseRabbitMQ()

	// declare queues and start background workers (consumers)
	if err := queue.DeclareQueues(); err != nil {
		log.Fatalf("Failed to declare queues: %v", err)
	}
	workers.StartAll()

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

	if err := r.Run(":" + envOrDefault("PORT", "8080")); err != nil {
		log.Fatal(err)
	}
}
