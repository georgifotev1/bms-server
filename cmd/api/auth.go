package main

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/georgifotev/bms/internal/store"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	ManagerRole = "manager"
	AdminRole   = "admin"
	StaffRole   = "staff"
)

type CreateUserPayload struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Role      string `json:"role"`
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

	if !IsValidRole(payload.Role) {
		app.badRequestResponse(w, r, errors.New("invalid user role"))
		return
	}

	userId, err := app.store.CreateUser(r.Context(), store.CreateUserParams{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     sql.NullString{String: payload.Email, Valid: payload.Email != ""},
		Password:  hashedPass,
		Role:      payload.Role,
	})
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	err = app.store.CreateNotification(r.Context(), store.CreateNotificationParams{
		Message: "New registration request",
		Roles:   []string{ManagerRole, AdminRole},
		UserID:  uuid.NullUUID{UUID: userId, Valid: true},
	})
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	err = writeJSON(w, http.StatusCreated, userId)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func IsValidRole(role string) bool {
	switch role {
	case ManagerRole, AdminRole, StaffRole:
		return true
	default:
		return false
	}
}
