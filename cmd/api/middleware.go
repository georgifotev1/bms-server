package main

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type contextKey string

const SessionContexKey = contextKey("session")

func (app *application) sessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(SessionCookie)
		if err != nil {
			app.unauthorizedErrorResponse(w, r, err)
			return
		}

		sessionId, err := uuid.Parse(cookie.Value)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		session, err := app.store.GetSessionById(r.Context(), sessionId)
		if err != nil {
			app.parseDBError(w, r, err)
			return
		}

		ctx := context.WithValue(r.Context(), SessionContexKey, session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
