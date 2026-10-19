package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type Artist struct {
	ID      int    `json:"id`
	Name    string `json:"name"`
	Country string `json:"country"`
}
type ErrorResponse struct {
	Error string
}

func registerRoutes(db *sql.DB) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /artists/{id}", func(w http.ResponseWriter, r *http.Request) {
		findArtistByIDHandler(db, w, r)
	})
	return mux
}
func writeJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}
func findArtistByID(ctx context.Context, db *sql.DB, id int) (Artist, error) {
	var artist Artist
	err := db.QueryRowContext(ctx, `
	SELECT artist_id,name,country
	FROM artists
	WHERE artist_id=?
	`, id).Scan(&artist.ID, &artist.Name, &artist.Country)
	if err != nil {
		return Artist{}, err
	}
	return artist, nil
}
func findArtistByIDHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	rawID := r.PathValue("id")
	id, err := strconv.Atoi(rawID)
	if id <= 0 || err != nil {
		writeJSONError(w, http.StatusBadRequest, "decode id failed")
		return
	}
	ctx := r.Context()
	artist, err := findArtistByID(ctx, db, id)
	if errors.Is(err, sql.ErrNoRows) {
		writeJSONError(w, http.StatusNotFound, "ErrNoRows")
	}
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "findArtistByID error")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&artist)
}
