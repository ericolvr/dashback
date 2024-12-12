package routes

import (
	"github.com/Alarmtekgit/websocket/internal/handlers"
	"github.com/Alarmtekgit/websocket/internal/middleware"
	"github.com/gin-gonic/gin"
)

func HistoryRoutes(
	router *gin.Engine,
	historyHandler *handlers.HistoryHandler,
	jwtSecret []byte,
) {
	routes := router.Group("/api/v1/history")
	routes.Use(middleware.AuthMiddleware(jwtSecret))
	{
		routes.POST("/", historyHandler.CreateHistory)
		routes.GET("/:equipmentid", historyHandler.FindHistoryByID)
		routes.GET("/filter/:equipmentid", historyHandler.GetHistoryByID)
	}
}
