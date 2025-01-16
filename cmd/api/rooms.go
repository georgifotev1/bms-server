package main

import (
	"net/http"

	"github.com/georgifotev/bms/internal/store"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (app *application) listRoomsHandler(w http.ResponseWriter, r *http.Request) {
	hotelId := chi.URLParam(r, "hotelId")

	id, err := uuid.Parse(hotelId)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	params := getQueryParams(r)

	rooms, err := app.store.ListRooms(r.Context(), store.ListRoomsParams{
		HotelID: id,
		Limit:   params.Limit,
		Offset:  params.Offset,
	})
	if err != nil {
		app.parseDBError(w, r, err)
		return
	}

	if err := writeJSON(w, http.StatusOK, rooms); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
