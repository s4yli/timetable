package events

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"middleware/example/internal/models"
	"net/http"
)

func Ctx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		eventId := chi.URLParam(r, "id")
		if eventId == "" {
			logrus.Error("missing event ID in URL")
			customError := &models.CustomError{
				Message: "missing event ID",
				Code:    http.StatusBadRequest,
			}
			w.WriteHeader(customError.Code)
			body, _ := json.Marshal(customError)
			_, _ = w.Write(body)
			return
		}

		ctx := context.WithValue(r.Context(), "eventId", eventId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
