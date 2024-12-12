package websocket

import (
	"sync"

	"github.com/gorilla/websocket"
)

var (
	clients = make(map[*websocket.Conn]bool)
	mutex   = sync.Mutex{}
)
