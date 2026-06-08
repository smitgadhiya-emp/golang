package routes

import (
	"gin-project/controllers"
	"gin-project/middleware"

	"github.com/gin-gonic/gin"
)

func authRoutes(r *gin.RouterGroup) {

	route := r.Group("/auth")

	route.POST("/signup", controllers.Singup)
	route.POST("/password-login", controllers.Login)

	// OTP
	route.POST("/send-otp", controllers.SendOTP)
	route.POST("/verify-otp", controllers.VerifyOTP)

	route.POST("/email/send-otp", controllers.SendOTP)
	route.POST("/email/verify-otp", controllers.VerifyOTP)

	route.GET("/oauth/google", controllers.GoogleOAuthStart)
	route.GET("/oauth/google/callback", controllers.GoogleOAuthCallback)

	route.POST("/change-password", middleware.ValidateToken, controllers.ChangePassword)
}
