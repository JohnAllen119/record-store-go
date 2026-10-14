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
	if dsn == "" {
		t.Fatalf("get dsn failed")
	}
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("database error: %v", err)
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		t.Fatalf("database error: %v", err)
	}
	return db
}

func TestRecordByIDInvalidID(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()
	mux := registerRoutes(db)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/records/abc", nil)
	mux.ServeHTTP(recorder, request)
	response := recorder.Result()
	defer response.Body.Close()
	if response.StatusCode != http.StatusBadRequest {
		t.Fatalf("expect:StatusBadRequest get:%d", response.StatusCode)
	}
}
func TestRecordByIDSucceed(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()
	mux := registerRoutes(db)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/records/1", nil)
	mux.ServeHTTP(recorder, request)
	response := recorder.Result()
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		t.Fatalf("expect:statusok,get:%v", response.StatusCode)
	}
	var record Record
	err := json.NewDecoder(response.Body).Decode(&record)
	if err != nil {
		t.Fatalf("database error: %v", err)
	}
	if record.ID != 1 || record.Title != "Abbey Road" || record.Artist != "The Beatles" {
		t.Fatalf("expect:id=1,title=abbeyroad,artist=the beatles")
	}
}
func TestRecordByIDNotFound(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()
	mux := registerRoutes(db)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/records/9999", nil)
	mux.ServeHTTP(recorder, request)
	response := recorder.Result()
	defer response.Body.Close()
	if response.StatusCode != http.StatusNotFound {
		t.Fatalf("expect:statusnotfound,get:%v", response.StatusCode)

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
		t.Fatalf("expected:statusinternalservererror,get:%d", response.StatusCode)
	}
	if response.Header.Get("Content-Type")!="application/json"{
		t.Fatalf("expect get spplication/json get %v",response.Header.Get("Content-Type"))
	}
}
