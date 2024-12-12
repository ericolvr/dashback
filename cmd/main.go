package main

import (
	"log"

	"github.com/Alarmtekgit/websocket/config"
	"github.com/Alarmtekgit/websocket/internal/consumers"
	"github.com/Alarmtekgit/websocket/internal/handlers"
	"github.com/Alarmtekgit/websocket/internal/operations"
	"github.com/Alarmtekgit/websocket/internal/repository"
	"github.com/Alarmtekgit/websocket/internal/routes"
	"github.com/Alarmtekgit/websocket/internal/service"
	"github.com/Alarmtekgit/websocket/internal/websocket"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	cfg := config.LoadConfig()

	client := config.GetMongoClient()
	db := client.Database(cfg.DatabaseName)

	jwtSecret := []byte(viper.GetString("JWT-SECRET"))

	// Get the collection
	collection := client.Database("dashboard").Collection("messages")

	messageRepo := repository.NewMessageRepository(db)
	messageService := service.NewMessageService(messageRepo)

	historyRepo := repository.NewHistoryRepository(db)
	historyService := service.NewHistoryService(historyRepo)

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://127.0.0.1:5173", "http://localhost:5173", "http://192.168.0.89:5173", "http://192.168.0.89"},
		AllowMethods:     []string{"GET", "PUT", "PATCH", "POST", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	routes.MessageRoutes(router, handlers.NewMessageHandler(messageService), jwtSecret)
	routes.HistoryRoutes(router, handlers.NewHistoryHandler(historyService), jwtSecret)
	routes.UserRoutes(router, handlers.NewUserHandler(userService), jwtSecret)

	router.GET("/ws", websocket.HandleWebSocket)
	router.GET("/sse", func(c *gin.Context) {
		operations.SSEHandler(c.Writer, c.Request)
	})

	go consumers.DatabaseConsumer(messageService, cfg.RabbitMQURL, collection)
	go consumers.WebSocketConsumer(cfg.RabbitMQURL)

	log.Printf("Service is running on port %s...", cfg.ServerPort)
	if err := router.Run(cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
