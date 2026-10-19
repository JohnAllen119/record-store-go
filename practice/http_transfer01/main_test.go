package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func OpenTestDB(t *testing.T) *sql.DB {
	dsn := os.Getenv("RECORD_STORE_HTTP_TEST_DSN")
	if dsn == "" {
		t.Fatalf("get dsn failed")
	}
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("mysql open failed %v", err)
	}
	err = db.Ping()
	if err != nil {
		t.Fatalf("mysql connect failed %v", err)
	}
	var dbName string
	err = db.QueryRow("SELECT DATABASE()").Scan(&dbName)
	if err != nil {
		t.Fatalf("select database failed")
	}
	if dbName != "record_store_http_test" {
		t.Fatalf("expect %s,got %s", "record_store_http_test", dbName)
	}
	return db
}
func TestDatabaseConnection(t *testing.T) {
	db := OpenTestDB(t)
	defer db.Close()
}
func TestGetArtistInvalidID(t *testing.T) {
	request := httptest.NewRequest("GET", "/artists/abc", nil)
	recorder := httptest.NewRecorder()
	mux := registerRoutes(nil)
	mux.ServeHTTP(recorder, request)
	response := recorder.Result()
	defer response.Body.Close()
	if response.StatusCode != http.StatusBadRequest {
		t.Fatalf("expect%v,got:%v", http.StatusBadRequest, response.StatusCode)
	}

}

func TestGetArtistNotFound(t *testing.T) {
	db := OpenTestDB(t)
	mux := registerRoutes(db)
	request := httptest.NewRequest("GET", "/artists/999999", nil)
	recorder := httptest.NewRecorder()
	mux.ServeHTTP(recorder, request)
	response := recorder.Result()
	defer response.Body.Close()
	if response.StatusCode != http.StatusNotFound {
		t.Fatalf("expect:%v,got:%v", http.StatusNotFound, response.StatusCode)
	}
}
func TestGetArtistSuccess(t *testing.T) {
	db := OpenTestDB(t)
	mux := registerRoutes(db)
	result, err := db.Exec(`
    INSERT INTO artists(name,country)
    VALUES(?,?)
`, "D'Angelo", "USA")
	if err != nil {
		t.Fatalf("insert failed:%v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("LastInsertId failed: %v", err)
	}
	t.Cleanup(func() {
		_, err := db.Exec(`
		DELETE FROM artists WHERE artist_id=?
		`, id)
		if err != nil {
			t.Errorf("cleanup failed: %v", err)
		}
	})
	path := fmt.Sprintf("/artists/%d", id)

	request := httptest.NewRequest("GET", path, nil)
	recorder := httptest.NewRecorder()
	mux.ServeHTTP(recorder, request)
	response := recorder.Result()
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		t.Fatalf("expect:statusok,get:%v", response.StatusCode)
	}
	var artist Artist
	err = json.NewDecoder(response.Body).Decode(&artist)
	if err != nil {
		t.Fatalf("decode failed %v", err)
	}
	if response.Header.Get("Content-Type") != "application/json" {
		t.Fatalf("expect application/json")
	}
	if artist.Name != "D'Angelo" || artist.Country != "USA" {
		t.Fatalf("expect %s,got %s", "D'Angelo", artist.Name)
	}
	if int64(artist.ID) != id {
		t.Fatalf("expect ID %d, got %d", id, artist.ID)
	}
}
