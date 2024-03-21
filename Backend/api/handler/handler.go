package handler

import (
	"context"
	"log"
	"net/http"

	"github.com/google/uuid"
)

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
