package routes

import (
	"gin-project/controllers"

	"github.com/gin-gonic/gin"
)

func userRoutes(r *gin.RouterGroup) {

	r.POST("/signup", controllers.Singup)
	r.POST("/login", controllers.Login)
}
