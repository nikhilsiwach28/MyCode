package handler

import (
	"context"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/nikhilsiwach28/MyCode.git/internal/user"
	"github.com/nikhilsiwach28/MyCode.git/models"
)

type userHandler struct {
	svc     user.Service
	sRouter *mux.Router
}

func NewUserHandler(svc user.Service) *userHandler {
	h := &userHandler{
		svc:     svc,
		sRouter: mux.NewRouter(),
	}

	return h.initRoutes()
}

func (h *userHandler) initRoutes() *userHandler {

	h.sRouter.HandleFunc("/user", serveHTTPWrapper(h.handleGet)).Methods("GET")
	h.sRouter.HandleFunc("/user", serveHTTPWrapper(h.handleCreate)).Methods("POST")
	// Add other routes as needed

	return h
}

func (h *userHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.sRouter.ServeHTTP(w, r)
}

func (h *userHandler) handleCreate(ctx context.Context, r *http.Request) apiResponse {
	var request models.CreateUserAPIRequest
	if err := request.Parse(r); err != nil {
		log.Println(err)
		return newAPIError(models.BadRequest.Add(err))
	}
	user := request.ToUser()
	user, err := h.svc.CreateUser(ctx, user)
	if err != models.NoError {
		return newAPIError(models.InternalError.Add(err))
	}

	return models.NewCreateUserAPIResponse(user)
}

func (h *userHandler) handleGet(ctx context.Context, r *http.Request) apiResponse {
	userID, err := uuid.Parse(r.URL.Query().Get("user_id"))
	if err != nil {
		return newAPIError(models.BadRequest.Add(err))
	}

	user, err := h.svc.GetUser(ctx, userID)
	if err != models.NoError {
		return newAPIError(models.InternalError.Add(err))
	}

	response := models.NewCreateUserAPIResponse(user)
	return response

}
