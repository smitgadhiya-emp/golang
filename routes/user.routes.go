package routes

import (
	"gin-project/controllers"
	"gin-project/middleware"

	"github.com/gin-gonic/gin"
)

func userRoutes(r *gin.RouterGroup) {
	// All user CRUD and profile endpoints require authentication
	r.Use(middleware.ValidateToken)

	r.GET("/users/me", controllers.GetMe)
	r.GET("/users", controllers.GetAllUsers)
	r.GET("/users/:id", controllers.GetUserByID)
	r.PUT("/users/:id", controllers.UpdateUser)
	r.DELETE("/users/:id", controllers.DeleteUser)
}
