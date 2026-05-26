package routes

import (
	"gin-project/controllers"
	"gin-project/middleware"

	"github.com/gin-gonic/gin"
)

func userRoutes(r *gin.RouterGroup) {

	r.POST("/signup", controllers.Singup)
	r.POST("/login", controllers.Login)

	r.GET("/oauth/google", controllers.GoogleOAuthStart)
	r.GET("/oauth/google/callback", controllers.GoogleOAuthCallback)

	r.POST("/change-password", middleware.ValidateToken, controllers.ChangePassword)
}
