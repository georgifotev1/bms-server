package main

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/georgifotev/bms/internal/store"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	ManagerRole   = "manager"
	AdminRole     = "admin"
	StaffRole     = "staff"
	PendingRole   = "pending"
	SessionCookie = "session_cookie"
)

type CreateUserPayload struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required"`
	Password  string `json:"password" validate:"required"`
}

func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateUserPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	err = app.store.CreateUser(r.Context(), store.CreateUserParams{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  hashedPass,
		Role:      PendingRole,
	})
	if err != nil {
		app.parseDBError(w, r, err)
		return
	}

	err = app.store.CreateNotification(r.Context(), store.CreateNotificationParams{
		Message: "New registration request",
		Roles:   []string{ManagerRole, AdminRole},
	})
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	err = writeJSON(w, http.StatusCreated, []string{})
	if err != nil {
		app.internalServerError(w, r, err)
	}
}

type CreateSessionPayload struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (app *application) createSessionHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateSessionPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user, err := app.store.GetUserByEmail(r.Context(), payload.Email)
	if err != nil {
		app.parseDBError(w, r, err)
		return
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(payload.Password))
	if err != nil {
		app.forbiddenResponse(w, r, err)
		return
	}

	var token uuid.UUID

	session, err := app.store.GetSessionByUserId(r.Context(), user.UserID)
	if errors.Is(err, sql.ErrNoRows) {
		sessionId, err := app.store.CreateSession(r.Context(), store.CreateSessionParams{
			SessionID: uuid.New(),
			UserID:    user.UserID,
			ExpiresAt: time.Now().UTC().Add(time.Hour),
		})
		if err != nil {
			app.parseDBError(w, r, err)
			return
		}
		token = sessionId
	} else if err != nil && !errors.Is(err, sql.ErrNoRows) {
		app.parseDBError(w, r, err)
		return
	}

	updatedSession, err := app.store.UpdateSession(r.Context(), store.UpdateSessionParams{
		SessionID: session.SessionID,
		ExpiresAt: time.Now().UTC().Add(time.Hour),
	})
	if err != nil {
		app.parseDBError(w, r, err)
		return
	}

	token = updatedSession

	cookie := &http.Cookie{
		Name:     SessionCookie,
		Value:    token.String(),
		Path:     "/",
		Expires:  time.Now().UTC().Add(time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, cookie)

	err = writeJSON(w, http.StatusCreated, []string{})
	if err != nil {
		app.internalServerError(w, r, err)
	}

}

func (app *application) clearSessionHandler(w http.ResponseWriter, r *http.Request) {
	ctxSession := r.Context().Value(SessionContexKey).(store.Session)

	_, err := app.store.UpdateSession(r.Context(), store.UpdateSessionParams{
		SessionID: ctxSession.SessionID,
		ExpiresAt: time.Unix(0, 0),
	})
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	expiredCookie := &http.Cookie{
		Name:     SessionCookie,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, expiredCookie)

	err = writeJSON(w, http.StatusOK, []string{})
	if err != nil {
		app.internalServerError(w, r, err)
	}
}

type UserResponse struct {
	UserID    string    `json:"user_id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	ctxSession := r.Context().Value(SessionContexKey).(store.Session)

	session, err := app.store.GetSessionById(r.Context(), ctxSession.SessionID)
	if err != nil {
		app.parseDBError(w, r, err)
	}

	user, err := app.store.GetUserById(r.Context(), session.UserID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	userResponse := UserResponse{
		UserID:    user.UserID.String(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}

	err = writeJSON(w, http.StatusOK, userResponse)
	if err != nil {
		app.internalServerError(w, r, err)
	}
}
