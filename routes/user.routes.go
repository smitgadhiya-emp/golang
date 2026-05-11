package routes

import (
	"gin-project/controllers"
	"gin-project/middleware"

	"github.com/gin-gonic/gin"
)

func userRoutes(r *gin.RouterGroup) {

	r.POST("/signup", controllers.Singup)
	r.POST("/login", controllers.Login)

	r.POST("/change-password", middleware.ValidateToken, controllers.ChangePassword)
}
