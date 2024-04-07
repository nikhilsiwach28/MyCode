package handler

import (
	"context"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/nikhilsiwach28/MyCode.git/internal/submission"
	"github.com/nikhilsiwach28/MyCode.git/models"
)

type submissionsHandler struct {
	svc     submission.Service
	sRouter *mux.Router
}

func NewSubmissionsHandler(svc submission.Service) *submissionsHandler {
	h := &submissionsHandler{
		svc:     svc,
		sRouter: mux.NewRouter(),
	}

	return h.initRoutes()
}

func (h *submissionsHandler) initRoutes() *submissionsHandler {

	h.sRouter.HandleFunc("/submission", serveHTTPWrapper(h.handleGet)).Methods("GET")
	h.sRouter.HandleFunc("/submission", serveHTTPWrapper(h.handleCreate)).Methods("POST")
	h.sRouter.HandleFunc("/submission/user", serveHTTPWrapper(h.handleGetUserSubmissions)).Methods("GET")
	// Add other routes as needed

	return h
}

func (h *submissionsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.sRouter.ServeHTTP(w, r)
}

func (h *submissionsHandler) handleCreate(ctx context.Context, r *http.Request) apiResponse {
	var request models.CreateSubmissionAPIRequest
	if err := request.Parse(r); err != nil {
		log.Println(err)
		return newAPIError(models.BadRequest.Add(err))
	}

	submission := request.ToSubmissions()
	submission, err := h.svc.CreateSubmission(ctx, submission, request.InputFile)
	if err != models.NoError {
		return newAPIError(models.InternalError.Add(err))
	}

	return models.NewCreateSubmissionAPIResponse(submission)
}

func (h *submissionsHandler) handleGet(ctx context.Context, r *http.Request) apiResponse {
	submissionID, err := uuid.Parse(r.URL.Query().Get("submission_id"))
	if err != nil {
		return newAPIError(models.BadRequest.Add(err))
	}

	submission, err := h.svc.GetSubmission(ctx, submissionID)
	if err != models.NoError {
		return newAPIError(models.InternalError.Add(err))
	}

	response := models.NewCreateSubmissionAPIResponse(submission)
	return response

}

func (h *submissionsHandler) handleGetUserSubmissions(ctx context.Context, r *http.Request) apiResponse {

	userId, err := uuid.Parse(r.URL.Query().Get("user_id"))
	if err != nil {
		return newAPIError(models.BadRequest.Add(err))
	}

	submissions, err := h.svc.GetUserSubmissions(ctx, userId)
	if err != models.NoError {
		return newAPIError(models.InternalError.Add(err))
	}

	response := models.NewUserSubmissionsAPIResponse(submissions)
	return response

}
