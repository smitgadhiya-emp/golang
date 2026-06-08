package routes

import (
	"gin-project/controllers"

	"github.com/gin-gonic/gin"
)

func emailRoutes(r *gin.RouterGroup) {

	route := r.Group("/email")

	route.POST("/welcome", controllers.SendWelcomeEmail)
	route.POST("/otp", controllers.SendOtpEmail)
}
