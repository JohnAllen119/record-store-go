package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

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
