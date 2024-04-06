package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/nikhilsiwach28/MyCode.git/api/handler"
	"github.com/nikhilsiwach28/MyCode.git/config"
	repo "github.com/nikhilsiwach28/MyCode.git/database"
	"github.com/nikhilsiwach28/MyCode.git/internal/submission"
	"github.com/nikhilsiwach28/MyCode.git/internal/user"
	"github.com/nikhilsiwach28/MyCode.git/models"
	"github.com/nikhilsiwach28/MyCode.git/queue"
	"github.com/nikhilsiwach28/MyCode.git/redis"
)

type APIServer struct {
	httpServer  *http.Server
	middlewares []mux.MiddlewareFunc
	router      *mux.Router
	rbac        map[http.Handler]models.AccessLevelModeEnum
}

func NewServer(cfg config.ServerConfig) *APIServer {
	return &APIServer{
		httpServer: &http.Server{
			Addr:         fmt.Sprintf("%s:%s", cfg.GetAddress(), cfg.GetPort()),
			WriteTimeout: time.Duration(cfg.GetWriteTimeout()) * time.Second,
			ReadTimeout:  time.Duration(cfg.GetReadTimeout()) * time.Second,
		},
		middlewares: []mux.MiddlewareFunc{},
		router:      mux.NewRouter(),
		rbac:        make(map[http.Handler]models.AccessLevelModeEnum),
	}
}

func (s *APIServer) add(path string, role models.AccessLevelModeEnum, handler http.Handler) {
	s.router.PathPrefix(path).Handler(handler)
	s.rbac[handler] = role
}

func (s *APIServer) initRoutesAndMiddleware() {

	// ADD Routes here
	connString := "host=localhost port=5432 user=username password=password dbname=database_name sslmode=disable"

	brokers := []string{"localhost:9092"}
	kafkaQueue := queue.InitQueue(brokers)
	redisClient := redis.NewRedisService("localhost:6379", "", 0)

	s.add("/submission", models.AccessLevelUser, handler.NewSubmissionsHandler(submission.New(repo.NewPostgres(connString))))
	s.add("/user", models.AccessLevelUser, handler.NewUserHandler(user.New(repo.NewPostgres(connString))))
	// s.add("/run", models.AccessLevelUser, s.handleRun())

	s.router.HandleFunc("/run", handler.NewRunHandler(kafkaQueue, redisClient)).Methods("POST", "GET")

	s.middlewares = []mux.MiddlewareFunc{
		mux.CORSMethodMiddleware(s.router),
		handler.NewReqIDMiddleware().Decorate,
		OptionMiddleware,
	}
	s.router.Use(s.middlewares...)
	s.httpServer.Handler = s.router
}

func (s *APIServer) handleRun() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("hello")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello from handleRun!"))
	}
}

func (s *APIServer) run() {
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("Error starting server:", err)
			os.Exit(1)
		}
	}()

	log.Println("[*] Server running .... ")

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	fmt.Println("Received signal:", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		fmt.Println("Error during server shutdown:", err)
	}
	fmt.Println("Server gracefully stopped")
}

func OptionMiddleware(next http.Handler) http.Handler {
	fmt.Print("hello from middleware")
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if req.Method == "OPTIONS" {
			w.WriteHeader(http.StatusAccepted)
			return
		}

		next.ServeHTTP(w, req)
	})
}

func StartHttpServer(cfg config.ServerConfig) {
	server := NewServer(cfg)
	server.initRoutesAndMiddleware()
	server.run()
}
