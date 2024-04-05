package main

import (
	"fmt"
	"log/slog"

	"github.com/joho/godotenv"
	"github.com/nikhilsiwach28/MyCode.git/api"
	"github.com/nikhilsiwach28/MyCode.git/config"
	// "github.com/nikhilsiwach28/MyCode.git/queue"
)

func main() {
	if err := godotenv.Load("local.env"); err != nil {
		slog.Warn("Error in loading env file, Generate .env file")
	}
	// queue.InitQueue()
	serverConf := config.NewServerConfig()
	fmt.Print(serverConf)
	api.StartHttpServer(serverConf)

}
