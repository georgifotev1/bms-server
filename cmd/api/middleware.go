package main

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/georgifotev/bms/internal/store"
	"github.com/google/uuid"
)

type contextKey string

const SessionContexKey = contextKey("session")

func (app *application) sessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := store.Session{
			SessionID: uuid.Nil,
			UserID:    uuid.Nil,
			ExpiresAt: time.Now().Add(-1 * time.Minute),
		}

		cookie, err := r.Cookie(SessionCookie)
		if err != nil {
			if err == http.ErrNoCookie {
				clearErr := app.store.ClearExpiredSessions(r.Context())
				if clearErr != nil {
					app.internalServerError(w, r, err)
					return
				}
			}
			ctx := context.WithValue(r.Context(), SessionContexKey, s)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		sessionId, err := uuid.Parse(cookie.Value)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		session, err := app.store.GetSessionById(r.Context(), sessionId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				ctx := context.WithValue(r.Context(), SessionContexKey, s)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			app.internalServerError(w, r, err)
			return
		}

		ctx := context.WithValue(r.Context(), SessionContexKey, session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
