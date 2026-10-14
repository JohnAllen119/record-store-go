package main

import (
	"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"strconv"
)

type Record struct {
	ID     int
	Title  string
	Artist string
}
type ErrorResponse struct {
	Error string
}

func registerRoutes(db *sql.DB) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /records/{id}", func(w http.ResponseWriter, r *http.Request) {
		recordByIDHandler(db, w, r)
	})
	return mux
}
func writeJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}
func findRecordByID(db *sql.DB, id int) (Record, error) {
	var record Record
	err := db.QueryRow(`
	SELECT record_id,title,artist
	FROM records
	WHERE record_id=?
	`, id).Scan(&record.ID, &record.Title, &record.Artist)
	if err != nil {
		return Record{}, err
	}
	return record, nil
}
func recordByIDHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	rawID := r.PathValue("id")
	id, err := strconv.Atoi(rawID)
	if id <= 0 || err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid id")
		return
	}
	record, err := findRecordByID(db, id)
	if err == sql.ErrNoRows {
		writeJSONError(w, http.StatusNotFound, "failed")
		return
	}
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "connect failed")
		return
	}
	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(&record)
	return
}
