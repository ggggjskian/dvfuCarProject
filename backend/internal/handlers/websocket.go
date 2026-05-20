package handlers

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var wsClients = make(map[string][]*websocket.Conn)
var wsMutex sync.Mutex
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func HandleWebSocket(c *gin.Context) {
	tripID := c.Param("id")
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	wsMutex.Lock()
	wsClients[tripID] = append(wsClients[tripID], conn)
	wsMutex.Unlock()

	defer func() {
		wsMutex.Lock()
		conns := wsClients[tripID]
		for i, c := range conns {
			if c == conn {
				wsClients[tripID] = append(conns[:i], conns[i+1:]...)
				break
			}
		}
		wsMutex.Unlock()
	}()

	for {
		var msg map[string]interface{}
		if err := conn.ReadJSON(&msg); err != nil {
			break
		}

		if msg["action"] == "location" {
			outMsg := map[string]interface{}{
				"type":      "driver_location",
				"lat":       msg["lat"],
				"lon":       msg["lon"],
				"driver_id": msg["driver_id"],
			}

			wsMutex.Lock()
			for _, client := range wsClients[tripID] {
				client.WriteJSON(outMsg)
			}
			wsMutex.Unlock()
		}
	}
}
