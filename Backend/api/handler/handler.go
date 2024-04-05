package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/nikhilsiwach28/MyCode.git/models"
)

type apiResponse interface {
	Write(http.ResponseWriter) error
}
type Handler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type ReqIDMiddleware struct {
	id uuid.UUID
}

func NewReqIDMiddleware() *ReqIDMiddleware {
	return &ReqIDMiddleware{
		id: uuid.New(),
	}
}

type Decorator interface {
	Decorate(handler Handler) Handler
}

func (middleware *ReqIDMiddleware) Decorate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ctx = context.WithValue(r.Context(), "ContextKeyRequestID", middleware.id.String())
		r = r.WithContext(ctx)
		log.Printf("Incomming request %s %s %s %s", r.Method, r.RequestURI, r.RemoteAddr, middleware.id.String())
		next.ServeHTTP(w, r)
		log.Printf("Finished handling http req. %s", middleware.id.String())
	})
}

func serveHTTPWrapper(f func(context.Context, *http.Request) apiResponse) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		var ctx = r.Context()
		response := f(ctx, r)

		// w.Header().Set("Access-Control-Allow-Origin", "*")
		if err := response.Write(w); err != nil {
			http.Error(w, "Failed to serialize response", http.StatusInternalServerError)
		}
	}
}

type apiError struct {
	status int
	body   models.AppError
}

func newAPIError(e models.AppError) *apiError {
	var err apiError
	err.body = e

	statusMapping := map[models.ErrorType]int{
		models.ErrorNone:          http.StatusOK,                  // 200 OK
		models.ErrorTimeout:       http.StatusRequestTimeout,      // 408 Request Timeout
		models.ErrorCanceled:      http.StatusRequestTimeout,      // 408 Request Timeout (or choose a different suitable code)
		models.ErrorExec:          http.StatusInternalServerError, // 500 Internal Server Error
		models.ErrorBadData:       http.StatusBadRequest,          // 400 Bad Request
		models.ErrorInternal:      http.StatusInternalServerError, // 500 Internal Server Error
		models.ErrorUnavailable:   http.StatusServiceUnavailable,  // 503 Service Unavailable
		models.ErrorNotFound:      http.StatusNotFound,            // 404 Not Found
		models.ErrorNotAcceptable: http.StatusNotAcceptable,       // 406 Not Acceptable
	}

	// Use a switch to set the err.status based on e.Type
	err.status = statusMapping[e.Type]

	return &err
}

func (e *apiError) Write(w http.ResponseWriter) error {
	// Implement serialization and writing logic for the User API response
	// Serialize the struct r and write it to the response writer
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.status)
	return json.NewEncoder(w).Encode(e.body)
}
