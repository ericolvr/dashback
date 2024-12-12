package routes

import (
	"github.com/Alarmtekgit/websocket/internal/handlers"
	"github.com/Alarmtekgit/websocket/internal/middleware"
	"github.com/gin-gonic/gin"
)

func MessageRoutes(
	router *gin.Engine,
	messageHandler *handlers.MessageHandler,
	jwtSecret []byte,

) {
	routes := router.Group("/api/v1/messages")
	routes.Use(middleware.AuthMiddleware(jwtSecret))
	{
		routes.POST("/", messageHandler.CreateMessage)
		routes.GET("/:id", messageHandler.GetMessageById)
		routes.PATCH("/:id", messageHandler.UpdateMessage)
		routes.DELETE("/delete/:id", messageHandler.DeleteMessage)
		routes.GET("/type/:type", messageHandler.GetMessagesByType)
		routes.GET("/monitored/:monitored", messageHandler.GetMonitoredMessages)
	}
}
