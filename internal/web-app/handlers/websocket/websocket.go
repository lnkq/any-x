package websocket

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"any-x/internal/models"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Storage interface {
	SavePost(ctx context.Context, post *models.Post) error
	GetPosts(ctx context.Context, limit int) ([]models.Post, error)
	Close() error
}

func New(log *slog.Logger, clients *map[*websocket.Conn]struct{}, mu *sync.RWMutex, s Storage) http.HandlerFunc {
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

		if s != nil {
			posts, err := s.GetPosts(context.Background(), 50)
			if err != nil {
				log.Error("Failed to get posts", slog.Any("error", err))
			} else {
				for i := len(posts) - 1; i >= 0; i-- {
					data, err := json.Marshal(posts[i])
					if err != nil {
						continue
					}
					if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
						log.Error("Write error sending history:", slog.Any("error", err))
						break
					}
				}
			}
		}

		for {
			messageType, message, err := conn.ReadMessage()
			if err != nil {
				log.Error("Read error:", slog.Any("error", err))
				break
			}
			log.Info("Received message:", slog.String("message", string(message)))

			post := &models.Post{
				Content:   string(message),
				CreatedAt: time.Now(),
			}

			if s != nil {
				if err := s.SavePost(context.Background(), post); err != nil {
					log.Error("Failed to save post", slog.Any("error", err))
				}
			}

			postData, err := json.Marshal(post)
			if err != nil {
				log.Error("Marshal error", slog.Any("error", err))
				continue
			}

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
				if err := client.WriteMessage(messageType, postData); err != nil {
					log.Error("Write error:", slog.Any("error", err))
				}
			}
		}
	}
}
