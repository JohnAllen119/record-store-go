package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Record struct {
	ID     int     `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}
type CreateRecordRequest struct {
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}
type UpdatePriceInput struct {
	Price float64 `json:"price"`
}
type ErrorResponse struct {
	Error string `json:"error"`
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
func createRecordHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var input CreateRecordRequest
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "JSON格式错误")
		return
	}

	input.Title = strings.TrimSpace(input.Title)
	input.Artist = strings.TrimSpace(input.Artist)

	if input.Title == "" || input.Artist == "" || input.Price <= 0 {
		writeJSONError(
			w,
			http.StatusBadRequest,
			"title、artist 不能为空，price 必须大于 0",
		)
		return
	}
	newID, err := addRecord(db, input.Title, input.Artist, input.Price)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "新增唱片失败")
		return
	}
	createdRecord := Record{
		ID:     int(newID),
		Title:  input.Title,
		Artist: input.Artist,
		Price:  input.Price,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdRecord)

}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "唱片服务运行中")
}
func recordsHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		var records []Record
		var err error
		query := r.URL.Query()
		var price float64
		if !query.Has("artist") && !query.Has("min_price") {
			records, err = getAllRecords(db)
		}
		if query.Has("min_price") {
			rawPrice := query.Get("min_price")
			var parseErr error
			price, parseErr = strconv.ParseFloat(rawPrice, 64)
			if parseErr != nil || price <= 0 || math.IsNaN(price) || math.IsInf(price, 0) {
				writeJSONError(w, http.StatusBadRequest, "invalid min_price")
				return
			}
			if !query.Has("artist") {
				records, err = getRecordsByMinPrice(db, price)
			}

		}
		if query.Has("artist") {
			artist := strings.TrimSpace(query.Get("artist"))
			if artist == "" {
				writeJSONError(w, http.StatusBadRequest, "invalid artist")
				return
			}
			if query.Has("min_price") {
				records, err = getRecordsByArtistAndMinPrice(db, artist, price)
			} else {
				records, err = getRecordsByArtist(db, artist)
			}

		}
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "查询唱片失败")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(records)

	case http.MethodPost:
		createRecordHandler(db, w, r)

	default:
		w.Header().Set("Allow", "GET, POST")
		writeJSONError(
			w,
			http.StatusMethodNotAllowed,
			"只允许 GET 或 POST 请求",
		)
	}
}
func recordByIDHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	rawID := r.PathValue("id")
	id, err := strconv.Atoi(rawID)
	if err != nil || id <= 0 {
		writeJSONError(w, http.StatusBadRequest, "invalid record id")
		return
	}
	record, err := findRecordByID(r.Context(), db, id)
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

func deleteRecordHandler(
	db *sql.DB,
	w http.ResponseWriter,
	r *http.Request,
) {
	rawID := r.PathValue("id")
	id, err := strconv.Atoi(rawID)
	if err != nil || id <= 0 {
		writeJSONError(w, http.StatusBadRequest, "invalid record id")
		return
	}
	affected, err := deleteRecord(db, id)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "数据库错误")
		return
	}
	if affected == 0 {
		writeJSONError(w, http.StatusNotFound, "没找到")
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
func updateRecordPriceHandler(
	db *sql.DB,
	w http.ResponseWriter,
	r *http.Request,
) {
	rawID := r.PathValue("id")
	id, err := strconv.Atoi(rawID)
	if err != nil || id <= 0 {
		writeJSONError(w, http.StatusBadRequest, "invalid record id")
		return
	}
	var input UpdatePriceInput
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if input.Price <= 0 {
		writeJSONError(w, http.StatusBadRequest, "价格必须大于0")
		return
	}
	affected, err := updateRecordPrice(db, id, input.Price)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "failed to update record")
		return
	}
	if affected == 0 {
		writeJSONError(w, http.StatusNotFound, "record not found")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(input)
}

func registerRoutes(db *sql.DB) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", helloHandler)
	mux.HandleFunc("/records", func(w http.ResponseWriter, r *http.Request) {
		recordsHandler(db, w, r)
	})

	mux.HandleFunc("GET /records/{id}", func(w http.ResponseWriter, r *http.Request) {
		recordByIDHandler(db, w, r)
	})
	mux.HandleFunc("PATCH /records/{id}", func(w http.ResponseWriter, r *http.Request) {
		updateRecordPriceHandler(db, w, r)
	})
	mux.HandleFunc("DELETE /records/{id}", func(w http.ResponseWriter, r *http.Request) {
		deleteRecordHandler(db, w, r)
	})
	return mux
}
func main() {
	dsn := os.Getenv("RECORD_STORE_DSN")
	if dsn == "" {
		fmt.Println("没有设置RECORD_STORE_DSN")
		return
	}
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("打开数据库失败：", err)
		return
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		fmt.Println("连接数据库失败：", err)
		return
	}
	fmt.Println("连接 MySQL 成功")
	mux := registerRoutes(db)

	fmt.Println("服务器启动：http://127.0.0.1:8080")

	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Println("服务器启动失败：", err)
	}
}
