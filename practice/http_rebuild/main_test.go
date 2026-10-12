package main

import (
	"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func openTestDB(t *testing.T) *sql.DB {
	dsn := os.Getenv("RECORD_STORE_DSN")
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		t.Fatalf("failed to ping database: %v", err)
	}
	return db
}
func TestRecordByIDInvalid(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/records/abc", nil)
	recorder := httptest.NewRecorder()
	mux := registerRoutes(nil)
	mux.ServeHTTP(recorder, request)
	response := recorder.Result()
	defer response.Body.Close()
	if response.StatusCode != http.StatusBadRequest {
		t.Fatalf("expect: 400,get: %v", response.StatusCode)
	}
	return
}

func TestRecordByIDNotFound(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()
	mux := registerRoutes(db)
	request := httptest.NewRequest(http.MethodGet, "/records/999999", nil)
	recorder := httptest.NewRecorder()
	mux.ServeHTTP(recorder, request)
	response := recorder.Result()
	defer response.Body.Close()
	if response.StatusCode != http.StatusNotFound {
		t.Fatalf("expected:404,get:%v", response.StatusCode)
	}
}

func TestRecordByIDSuccess(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()
	mux := registerRoutes(db)
	request := httptest.NewRequest(http.MethodGet, "/records/1", nil)
	recorder := httptest.NewRecorder()
	mux.ServeHTTP(recorder, request)
	response := recorder.Result()
	defer response.Body.Close()
	var record Record
	json.NewDecoder(response.Body).Decode(&record)
	if record.ID != 1 || record.Title != "Abbey Road" || record.Artist != "The Beatles" {
		t.Fatalf("expect : id=1 get:%d title :Abbey Road  get : %s artist : The Beatles get:%s", record.ID, record.Title, record.Artist)
	}
}

func TestRecordByIDDatabaseError(t *testing.T) {
	db := openTestDB(t)
	db.Close()
	mux := registerRoutes(db)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/records/1", nil)
	mux.ServeHTTP(recorder, request)
	response := recorder.Result()
	defer response.Body.Close()
	if response.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expect:StatusInternalServerError get:%d", response.StatusCode)
	}
	if a := response.Header.Get("Content-Type"); a != "application/json" {
		t.Fatalf("expect:content-type")
	}

}
