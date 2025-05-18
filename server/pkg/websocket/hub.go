package websocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Hub struct {
	clients   map[*websocket.Conn]bool
	broadcast chan string
	lock      sync.Mutex
}

var hub = &Hub{
	clients:   make(map[*websocket.Conn]bool),
	broadcast: make(chan string),
}

func HubInstance() *Hub {
	return hub
}

func (h *Hub) Run() {
	for {
		msg := <-h.broadcast

		h.lock.Lock()
		for conn := range h.clients {
			if conn == nil {
				delete(h.clients, conn)
				continue
			}
			err := conn.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				fmt.Println("Error writing to WebSocket:", err)
				conn.Close()
				delete(h.clients, conn)
			}
		}
		h.lock.Unlock()
	}
}

func (h *Hub) HandleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("WebSocket upgrade failed:", err)
		return
	}

	h.lock.Lock()
	h.clients[conn] = true
	h.lock.Unlock()

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			h.lock.Lock()
			delete(h.clients, conn)
			h.lock.Unlock()
			conn.Close()
			break
		}
	}
}

func (h *Hub) SendUpdate(data string) {
	h.broadcast <- data
}
