package main

import (
	"net/http"

	"github.com/georgifotev/bms/internal/store"
	"golang.org/x/crypto/bcrypt"
)

const (
	ManagerRole = "manager"
	AdminRole   = "admin"
	StaffRole   = "staff"
	PendingRole = "pending"
)

type CreateUserPayload struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required"`
	Password  string `json:"password" validate:"required"`
	// Role      string `json:"role" validate:"manager|admin|staff|pending"`
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

	userId, err := app.store.CreateUser(r.Context(), store.CreateUserParams{
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

	err = writeJSON(w, http.StatusCreated, userId)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func IsValidRole(role string) bool {
	switch role {
	case ManagerRole, AdminRole, StaffRole, PendingRole:
		return true
	default:
		return false
	}
}
