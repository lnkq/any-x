package main

import (
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"

	websocketHandler "any-x/internal/web-app/handlers/websocket"

	"github.com/MatusOllah/slogcolor"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]struct{})
var clientsMu sync.RWMutex

func main() {
	log := setupLogger()

	router := chi.NewRouter()
	router.HandleFunc("/ws", websocketHandler.New(log, &clients, &clientsMu))
	router.Handle("/", http.FileServer(http.Dir("./static")))

	log.Info("Server starting on :8080")

	go health(log)

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Error("ListenAndServe error:", slog.Any("error", err))
	}
}

func health(log *slog.Logger) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		clientsMu.RLock()
		n := len(clients)
		clientsMu.RUnlock()
		log.Info("Active connections", slog.Int("count", n))
	}
}

func setupLogger() *slog.Logger {
	return slog.New(slogcolor.NewHandler(os.Stderr, slogcolor.DefaultOptions))
}
