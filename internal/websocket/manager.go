package websocket

import (
	"encoding/json"
	"log"

	"github.com/Alarmtekgit/websocket/internal/domain"
	"github.com/gorilla/websocket"
)

func BroadcastMessage(message domain.Message) {
	// fmt.Println("Broadcasting message:", message)
	mutex.Lock()
	defer mutex.Unlock()

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		return
	}

	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, jsonMessage)
		if err != nil {
			log.Println("Error sending message to client:", err)
			client.Close()
			delete(clients, client)
		}
	}
}
