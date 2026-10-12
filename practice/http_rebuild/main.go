package main

import (
	"database/sql"
	"encoding/json"
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

func writeJSONError(
	w http.ResponseWriter,
	status int,
	message string,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}
func recordByIDHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	rawID := r.PathValue("id")
	id, err := strconv.Atoi(rawID)
	if err != nil || id <= 0 {
		writeJSONError(w, http.StatusBadRequest, "invalid id")
		return
	}

	record, err := findRecordByID(db, id)
	if err == sql.ErrNoRows {
		writeJSONError(w, http.StatusNotFound, "record not found")
		return
	}
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "failed to query record")
		return

	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(record)
}
func findRecordByID(db *sql.DB, recordID int) (Record, error) {
	var record Record
	err := db.QueryRow(`
	SELECT record_id,title,artist
	FROM records
	WHERE record_id = ?
	ORDER BY record_id
	`, recordID).Scan(&record.ID, &record.Title, &record.Artist)

	if err != nil {
		return Record{}, err
	}
	return record, nil
}
func registerRoutes(db *sql.DB) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /records/{id}", func(w http.ResponseWriter, r *http.Request) {
		recordByIDHandler(db, w, r)
	})
	return mux
}
func main() {

}
