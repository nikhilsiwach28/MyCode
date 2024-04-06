package handler

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/nikhilsiwach28/MyCode.git/internal/run"
	"github.com/nikhilsiwach28/MyCode.git/queue"
	"github.com/nikhilsiwach28/MyCode.git/redis"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewRunHandler(queueService queue.QueueService, redisService *redis.RedisService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Upgrade HTTP connection to WebSocket
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, "Failed to upgrade to WebSocket", http.StatusInternalServerError)
			return
		}

		// Handle run logic
		run.HandleRun(ws, queueService, redisService)
	}
}
