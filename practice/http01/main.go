package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
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

func addRecord(db *sql.DB,
	title string,
	artist string,
	price float64) (int64, error) {
	result, err := db.Exec(`
	INSERT INTO records(title,artist,price)
	VALUES(?,?,?)
		`, title, artist, price)
	if err != nil {
		return 0, err
	}
	newID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return newID, nil
}
func getAllRecords(db *sql.DB) ([]Record, error) {
	var records []Record
	rows, err := db.Query(`
	SELECT record_id,title,artist,price
	FROM records
	ORDER BY record_id
	`)
	if err != nil {
		return []Record{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var record Record
		err := rows.Scan(&record.ID, &record.Title, &record.Artist, &record.Price)
		if err != nil {
			return []Record{}, err
		}
		records = append(records, record)
	}
	if err := rows.Err(); err != nil {
		return []Record{}, err
	}
	return records, nil

}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "唱片服务运行中")
}
func recordsHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		records, err := getAllRecords(db)
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
	SELECT record_id,title,artist,price
	FROM records
	WHERE record_id = ?;
	`, recordID).Scan(&record.ID, &record.Title, &record.Artist, &record.Price)
	if err != nil {
		return Record{}, err
	}
	return record, nil

}

func deleteRecord(db *sql.DB, recordID int) (int64, error) {
	result, err := db.Exec(`
	DELETE FROM records
	WHERE record_id = ?;
	`, recordID)
	if err != nil {
		return 0, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return affected, nil

}
func updateRecordPrice(
	db *sql.DB,
	recordID int,
	newPrice float64,
) (int64, error) {
	result, err := db.Exec(`
	UPDATE records
	SET price=?
	WHERE record_id=?;
	`, newPrice, recordID)
	if err != nil {
		return 0, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return affected, nil
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
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/records", func(w http.ResponseWriter, r *http.Request) {
		recordsHandler(db, w, r)
	})

	http.HandleFunc("GET /records/{id}", func(w http.ResponseWriter, r *http.Request) {
		recordByIDHandler(db, w, r)
	})
	http.HandleFunc("PATCH /records/{id}", func(w http.ResponseWriter, r *http.Request) {
		updateRecordPriceHandler(db, w, r)
	})
	http.HandleFunc("DELETE /records/{id}", func(w http.ResponseWriter, r *http.Request) {
		deleteRecordHandler(db, w, r)
	})
	fmt.Println("服务器启动：http://127.0.0.1:8080")

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("服务器启动失败：", err)
	}
}
