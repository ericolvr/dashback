package routes

import (
	"github.com/Alarmtekgit/websocket/internal/handlers"
	"github.com/Alarmtekgit/websocket/internal/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(
	router *gin.Engine,
	userHandler *handlers.UserHandler,
	jwtSecret []byte,
) {
	routes := router.Group("/api/v1/users")
	routes.POST("/authenticate", userHandler.Authenticate)
	routes.Use(middleware.AuthMiddleware(jwtSecret))
	{
		routes.POST("/", userHandler.CreateUser)
		routes.GET("/:id", userHandler.GetUserById)
		routes.GET("/mobile/:mobile", userHandler.GetUserByMobile)
		routes.DELETE("/delete/:id", userHandler.DeleteUser)
		routes.GET("/", userHandler.GetAllUsers)
	}
}
