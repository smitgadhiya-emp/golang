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

	// Password change (requires authentication)
	route.POST("/change-password", middleware.ValidateToken, controllers.ChangePassword)
	
	// forgot password
	route.POST("/forgot-password", controllers.ForgotPassword)
	route.POST("/forgot-password/verify-token", controllers.VerifyResetToken)
	route.POST("/forgot-password/reset", controllers.ResetPassword)

	// Google OAuth
	
	route.GET("/oauth/google", controllers.GoogleOAuthStart)
	route.GET("/oauth/google/callback", controllers.GoogleOAuthCallback)
}
