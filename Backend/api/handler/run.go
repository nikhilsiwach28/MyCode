package handler

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/nikhilsiwach28/MyCode.git/internal/run"
	"github.com/nikhilsiwach28/MyCode.git/queue"
	"github.com/nikhilsiwach28/MyCode.git/redis"
)

func NewRunHandler(webSocket websocket.Upgrader, queueService queue.QueueService, redisService *redis.RedisService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Upgrade HTTP connection to WebSocket
		fmt.Println("NewRunHandler")
		ws, err := webSocket.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, "Failed to upgrade to WebSocket", http.StatusInternalServerError)
			return
		}

		// Handle run logic
		run.HandleRun(ws, queueService, redisService)
	}
}
