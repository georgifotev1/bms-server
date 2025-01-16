package main

import (
	"net/http"
	"strconv"
)

type Params struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func getQueryParams(r *http.Request) Params {
	limit := 10
	offset := 0

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil {
			return Params{Limit: int32(limit), Offset: int32(offset)}
		}
		limit = parsedLimit
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		parsedOffset, err := strconv.Atoi(offsetStr)
		if err != nil {
			return Params{Limit: int32(limit), Offset: int32(offset)}
		}
		offset = parsedOffset
	}
	return Params{Limit: int32(limit), Offset: int32(offset)}
}
