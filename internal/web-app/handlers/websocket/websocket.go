package websocket

import (
	"log/slog"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func New(log *slog.Logger, clients *map[*websocket.Conn]struct{}, mu *sync.RWMutex) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Error("Upgrade error:", slog.Any("error", err))
			return
		}
		defer conn.Close()

		mu.Lock()
		(*clients)[conn] = struct{}{}
		mu.Unlock()
		defer func() {
			mu.Lock()
			delete(*clients, conn)
			mu.Unlock()
		}()

		log.Info("New client connected")

		for {
			messageType, message, err := conn.ReadMessage()
			if err != nil {
				log.Error("Read error:", slog.Any("error", err))
				break
			}
			log.Info("Received message:", slog.String("message", string(message)))

			mu.RLock()
			targets := make([]*websocket.Conn, 0, len(*clients))
			for client := range *clients {
				if client == conn {
					continue
				}
				targets = append(targets, client)
			}
			mu.RUnlock()

			for _, client := range targets {
				if err := client.WriteMessage(messageType, message); err != nil {
					log.Error("Write error:", slog.Any("error", err))
				}
			}
		}
	}
}
