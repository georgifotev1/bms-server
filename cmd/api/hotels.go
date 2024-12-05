package main

import (
	"net/http"

	"github.com/georgifotev/bms/internal/store"
)

func (app *application) listHotelsHandler(w http.ResponseWriter, r *http.Request) {
	defaultParams := store.ListHotelsParams{
		Limit:  10,
		Offset: 0,
	}

	hotels, err := app.store.ListHotels(r.Context(), defaultParams)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := writeJSON(w, http.StatusOK, hotels); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
