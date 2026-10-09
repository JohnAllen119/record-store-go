package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func openTestDB(t *testing.T) *sql.DB {
	dsn := os.Getenv("RECORD_STORE_DSN")
	if dsn == "" {
		t.Fatal("RECORD_STORE_DSN is not set")
	}
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	if err := db.Ping(); err != nil {
		db.Close()
		t.Fatalf("failed to connect database: %v", err)
	}
	return db
}
func TestWriteJSONError(t *testing.T) {
	//创建recorder
	recorder := httptest.NewRecorder()
	//把recorder当作w传给writeJSONError
	writeJSONError(recorder, http.StatusBadRequest, "invalid input")
	//recorder保存响应
	response := recorder.Result()
	defer response.Body.Close()
	//下面就是检查状态码，响应体，响应头
	var body ErrorResponse
	err := json.NewDecoder(response.Body).Decode(&body)
	if err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	wantError := "invalid input"
	if body.Error != wantError {
		t.Fatalf("expected error message %q, got %q", wantError, body.Error)
	}
	contentType := response.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %q", contentType)
	}
	if response.StatusCode != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			response.StatusCode,
		)
	}
}
func TestRecordsHandlerMethodNotAllowed(t *testing.T) {
	request := httptest.NewRequest(http.MethodPut, "/records", nil)
	recorder := httptest.NewRecorder()
	recordsHandler(nil, recorder, request)
	response := recorder.Result()
	defer response.Body.Close()
	if response.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusMethodNotAllowed,
			response.StatusCode,
		)
	}
	wantAllow := "GET, POST"
	allow := response.Header.Get("Allow")
	if allow != wantAllow {
		t.Fatalf("expected Allow %q, got %q", wantAllow, allow)
	}
	var body ErrorResponse
	err := json.NewDecoder(response.Body).Decode(&body)
	if err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	wantError := "只允许 GET 或 POST 请求"
	if body.Error != wantError {
		t.Fatalf("expected error %q, got %q", wantError, body.Error)
	}
}
func TestUpdateRecordPriceHandlerDatabaseError(t *testing.T) {
	db, err := sql.Open("mysql", "")
	if err != nil {
		t.Fatalf("failed to create database handle: %v", err)
	}
	err = db.Close()
	if err != nil {
		t.Fatalf("failed to close database handle: %v", err)
	}
	request := httptest.NewRequest(
		http.MethodPatch,
		"/records/1",
		strings.NewReader(`{"price":99.99}`),
	)
	request.SetPathValue("id", "1")
	recorder := httptest.NewRecorder()
	updateRecordPriceHandler(db, recorder, request)
	response := recorder.Result()
	defer response.Body.Close()
	if response.StatusCode != http.StatusInternalServerError {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusInternalServerError,
			response.StatusCode,
		)
	}
}

func TestUpdateRecordPriceHandlerInvalidID(t *testing.T) {
	request := httptest.NewRequest(http.MethodPatch, "/records/abc", nil)
	request.SetPathValue("id", "abc")
	recorder := httptest.NewRecorder()
	updateRecordPriceHandler(nil, recorder, request)
	response := recorder.Result()
	defer response.Body.Close()
	if response.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status %d,get %d", http.StatusBadRequest, response.StatusCode)
	}

	var body ErrorResponse
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if body.Error != "invalid record id" {
		t.Fatalf("expected error %q, got %q", "invalid record id", body.Error)
	}
}
func TestRecordByIDHandlerInvalidID(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/records/abc", nil)
	request.SetPathValue("id", "abc")
	recorder := httptest.NewRecorder()
	recordByIDHandler(nil, recorder, request)
	response := recorder.Result()
	defer response.Body.Close()
	if response.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status %d,get %d", http.StatusBadRequest, response.StatusCode)
	}
	var body ErrorResponse
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	if body.Error != "invalid record id" {
		t.Fatalf("expected error %q, got %q", "invalid record id", body.Error)
	}

}

func TestRegisterRoutesInvalidRecordID(t *testing.T) {
	mux := registerRoutes(nil)
	request := httptest.NewRequest(http.MethodGet, "/records/abc", nil)
	recorder := httptest.NewRecorder()
	mux.ServeHTTP(recorder, request)
	response := recorder.Result()
	defer response.Body.Close()
	if response.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, response.StatusCode)
	}
	var body ErrorResponse
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	if body.Error != "invalid record id" {
		t.Fatalf("expected error %q, got %q", "invalid record id", body.Error)
	}
}
func TestRegisterRoutesInvalidArtist(t *testing.T) {
	mux := registerRoutes(nil)
	request := httptest.NewRequest(http.MethodGet, "/records?artist=%20%20", nil)
	recorder := httptest.NewRecorder()
	mux.ServeHTTP(recorder, request)
	response := recorder.Result()
	defer response.Body.Close()
	if response.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected:%d,get:%d", http.StatusBadRequest, response.StatusCode)
	}
	var body ErrorResponse
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	if body.Error != "invalid artist" {
		t.Fatalf("expected error %q, got %q", "invalid artist", body.Error)
	}
}
func TestRegisterRoutesInvalidMin_Price(t *testing.T) {
	mux := registerRoutes(nil)

	for _, raw := range []string{"abc", "NaN", "Inf"} {
		request := httptest.NewRequest(
			http.MethodGet, "/records?min_price="+raw, nil,
		)
		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, request)

		response := recorder.Result()
		defer response.Body.Close()

		if response.StatusCode != http.StatusBadRequest {
			t.Fatalf("min_price=%q: expected 400, got %d", raw, response.StatusCode)
		}

		var body ErrorResponse
		if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
			t.Fatalf("min_price=%q: decode response: %v", raw, err)
		}
		if body.Error != "invalid min_price" {
			t.Fatalf("min_price=%q: expected invalid min_price, got %q", raw, body.Error)
		}
	}
}

func TestRegisterRoutesQueryRecordsSuccess(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()
	newID, err := addRecord(
		db,
		"HTTP Test Album",
		"HTTP Test Artist",
		123.45,
	)
	if err != nil {
		t.Fatalf("failed to add test record: %v", err)
	}
	defer func() {
		if _, err := deleteRecord(db, int(newID)); err != nil {
			t.Errorf("failed to delete test record: %v", err)
		}
	}()
	mux := registerRoutes(db)
	request := httptest.NewRequest(
		http.MethodGet,
		"/records?artist=HTTP%20Test%20Artist",
		nil,
	)
	recorder := httptest.NewRecorder()
	mux.ServeHTTP(recorder, request)
	response := recorder.Result()
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			response.StatusCode,
		)
	}
	var records []Record
	if err := json.NewDecoder(response.Body).Decode(&records); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}
	record := records[0]

	if record.ID != int(newID) {
		t.Fatalf("expected record id %d, got %d", newID, record.ID)
	}

	if record.Title != "HTTP Test Album" {
		t.Fatalf("expected title %q, got %q", "HTTP Test Album", record.Title)
	}

	if record.Artist != "HTTP Test Artist" {
		t.Fatalf("expected artist %q, got %q", "HTTP Test Artist", record.Artist)
	}

	if record.Price != 123.45 {
		t.Fatalf("expected price %.2f, got %.2f", 123.45, record.Price)
	}

}

func TestRegisterRoutesQueryRecordsEmpty(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()
	mux := registerRoutes(db)
	request := httptest.NewRequest(
		http.MethodGet,
		"/records?artist=__HTTP_TEST_NON_EXISTENT_ARTIST__",
		nil,
	)
	recorder := httptest.NewRecorder()
	mux.ServeHTTP(recorder, request)
	response := recorder.Result()
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			response.StatusCode,
		)
	}
	var records []Record
	if err := json.NewDecoder(response.Body).Decode(&records); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	if len(records) != 0 {
		t.Fatalf("expect %d,got %d", 0, len(records))
	}
}

func TestFindRecordByIDCanceledContext(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := findRecordByID(ctx, db, 1)

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected: %v, got: %v", context.Canceled, err)
	}

}

func TestRecordByIDHandlerDeadlineExceeded(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()
	mux := registerRoutes(db)
	ctx, cancel := context.WithTimeout(context.Background(), 0*time.Second)
	defer cancel()
	request := httptest.NewRequest(http.MethodGet, "/records/10", nil)
	request = request.WithContext(ctx)
	recorder := httptest.NewRecorder()
	mux.ServeHTTP(recorder, request)
	response := recorder.Result()
	defer response.Body.Close()
	if response.StatusCode != http.StatusServiceUnavailable {
		t.Fatalf("expected: %v, got: %v", http.StatusServiceUnavailable, response.StatusCode)
	}
}

func TestRecordByIDHandlerCanceledContext(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()
	mux := registerRoutes(db)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	request := httptest.NewRequest(http.MethodGet, "/records/10", nil)
	request = request.WithContext(ctx)
	recorder := httptest.NewRecorder()
	mux.ServeHTTP(recorder, request)
	response := recorder.Result()
	defer response.Body.Close()
	if recorder.Body.Len() != 0 {
		t.Fatalf("expected len=0")
	}
}
