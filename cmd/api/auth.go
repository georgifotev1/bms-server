package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/georgifotev/bms/internal/store"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	ManagerRole = "manager"
	AdminRole   = "admin"
	StaffRole   = "staff"
	PendingRole = "pending"
)

type UserWithToken struct {
	UserID    uuid.UUID `json:"user_id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	Token     string    `json:"token"`
}

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

	user, err := app.store.CreateUser(r.Context(), store.CreateUserParams{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  hashedPass,
		Role:      PendingRole,
	})
	if err != nil {
		app.internalServerError(w, r, err)
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

	userWithoutToken := UserWithToken{
		UserID:    user.UserID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		Token:     "",
	}

	err = writeJSON(w, http.StatusCreated, userWithoutToken)
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
		app.logger.Error(err) // todo
		return
	}

	if user.Role == PendingRole {
		app.forbiddenResponse(w, r, errors.New("pending"))
		return
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(payload.Password))
	if err != nil {
		app.unauthorizedErrorResponse(w, r, err)
		return
	}

	token, err := app.store.CreateSession(r.Context(), store.CreateSessionParams{
		SessionID: uuid.New(),
		UserID:    user.UserID,
		ExpiresAt: time.Now().Add(time.Hour),
	})
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	userWithToken := UserWithToken{
		UserID:    user.UserID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		Token:     token.String(),
	}

	err = writeJSON(w, http.StatusCreated, userWithToken)
	if err != nil {
		app.internalServerError(w, r, err)
	}
}
