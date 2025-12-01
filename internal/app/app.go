package app

import (
	"log/slog"
	"net/http"
	"os"
	"sync"

	"any-x/internal/config"
	"any-x/internal/storage/sqlite"
	websocketHandler "any-x/internal/web-app/handlers/websocket"

	"github.com/MatusOllah/slogcolor"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

type App struct {
	log    *slog.Logger
	cfg    *config.Config
	server *http.Server
}

func New(cfg *config.Config) *App {
	log := setupLogger(cfg.Env)

	var store *sqlite.Storage
	var err error

	if cfg.StoragePath != "" {
		store, err = sqlite.New(cfg.StoragePath)
		if err != nil {
			log.Error("failed to init storage", slog.Any("err", err))
			os.Exit(1)
		}
	} else {
		log.Warn("No storage path provided, using nil storage")
	}

	router := chi.NewRouter()

	clients := make(map[*websocket.Conn]struct{})
	var clientsMu sync.RWMutex

	router.HandleFunc("/ws", websocketHandler.New(log, &clients, &clientsMu, store))
	router.Handle("/", http.FileServer(http.Dir("./static")))

	server := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	return &App{
		log:    log,
		cfg:    cfg,
		server: server,
	}
}

func (a *App) Run() error {
	a.log.Info("Server starting", slog.String("addr", a.cfg.HTTPServer.Address))
	return a.server.ListenAndServe()
}

func setupLogger(env string) *slog.Logger {
	return slog.New(slogcolor.NewHandler(os.Stderr, slogcolor.DefaultOptions))
}
