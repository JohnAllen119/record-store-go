package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"os"
)

type Record struct {
	ID     int
	Title  string
	Artist string
	Price  float64
}
type RecordStats struct {
	Count        int
	AveragePrice float64
	LowestPrice  float64
	HighestPrice float64
}
type ArtistStats struct {
	Artist       string
	Count        int
	AveragePrice float64
}

func getArtistStats(db *sql.DB) ([]ArtistStats, error) {
	var artistStates []ArtistStats
	rows, err := db.Query(`
	SELECT
		artist,
		COUNT(*),
		AVG(price)
	FROM records
	GROUP BY artist
	ORDER BY AVG(price) DESC
	`)
	if err != nil {
		return []ArtistStats{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var artistState ArtistStats
		err := rows.Scan(&artistState.Artist, &artistState.Count, &artistState.AveragePrice)
		if err != nil {
			return []ArtistStats{}, err
		}
		artistStates = append(artistStates, artistState)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return artistStates, nil
}

func getRecordStats(db *sql.DB) (RecordStats, error) {
	var recordState RecordStats
	err := db.QueryRow(`
	SELECT
		COUNT(*),
		AVG(price),
		MIN(price),
		MAX(price)
	FROM records;
	`).Scan(&recordState.Count, &recordState.AveragePrice, &recordState.LowestPrice, &recordState.HighestPrice)

	if err != nil {
		return RecordStats{}, err
	}
	return recordState, nil
}
func deleteRecord(db *sql.DB, recordID int) (int64, error) {
	result, err := db.Exec(`
	DELETE FROM records
	WHERE record_id=?
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
		WHERE record_id=?
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
func addRecord(
	db *sql.DB,
	title string,
	artist string,
	price float64,
) (int64, error) {
	result, err := db.Exec(`
		INSERT INTO records (title,artist,price)
		VALUES (?,?,?)
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
func findRecordsByArtistAndMinPrice(
	db *sql.DB,
	artistName string,
	minPrice float64,
) ([]Record, error) {
	var records []Record

	rows, err := db.Query(`
	SELECT record_id,title,artist,price
	FROM records
	WHERE artist=? AND price>=?
	ORDER BY record_id
	`, artistName, minPrice)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var record Record

		err = rows.Scan(
			&record.ID,
			&record.Title,
			&record.Artist,
			&record.Price,
		)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return records, nil
}
func findRecordByID(db *sql.DB, recordID int) (Record, error) {
	var record Record
	err := db.QueryRow(`
	SELECT record_id,title,artist,price
	FROM records
	WHERE record_id = ?
	`, recordID).Scan(&record.ID, &record.Title, &record.Artist, &record.Price)
	return record, err
}
func findLatestRecords(db *sql.DB, limit int) ([]Record, error) {
	var records []Record
	rows, err := db.Query(`
	SELECT record_id,title,artist,price
	FROM records
	ORDER BY record_id DESC
	LIMIT ?
	`, limit)

	if err != nil {
		return records, err
	}
	defer rows.Close()
	for rows.Next() {
		var record Record
		err := rows.Scan(&record.ID, &record.Title, &record.Artist, &record.Price)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return records, nil
}

func main() {
	dsn := os.Getenv("RECORD_STORE_DSN")
	if dsn == "" {
		fmt.Println("未设置环境变量 RECORD_STORE_DSN")
		return
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("创建数据库连接失败：", err)
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Println("连接 MySQL 失败：", err)
		return
	}

	fmt.Println("连接 MySQL 成功")

	minPrice := 140.00
	artistName := "The Beatles"
	records, err := findRecordsByArtistAndMinPrice(
		db,
		artistName,
		minPrice,
	)

	if err != nil {
		fmt.Println("查询数据失败：", err)
		return
	}
	for _, record := range records {
		fmt.Printf(
			"ID: %d | %s | %s | %.2f 元\n",
			record.ID,
			record.Title,
			record.Artist,
			record.Price)
	}
	fmt.Printf("本次共读取%d张唱片\n", len(records))
	record, err := findRecordByID(db, 2)

	if err == sql.ErrNoRows {
		fmt.Println("没有找到这张唱片", err)
		return
	}

	if err != nil {
		fmt.Println("这张唱片查询失败", err)
		return
	}

	lastRecords, err := findLatestRecords(db, 3)
	if err != nil {
		fmt.Println("查询最新唱片失败:", err)
		return
	}
	for _, record := range lastRecords {
		fmt.Printf(
			"ID: %d | %s | %s | %.2f 元\n",
			record.ID,
			record.Title,
			record.Artist,
			record.Price,
		)
	}
	fmt.Printf("找到唱片:%d | %s | %s | %.2f 元\n", record.ID, record.Title, record.Artist, record.Price)
	stats, err := getRecordStats(db)
	if err != nil {
		fmt.Println("统计失败：", err)
		return
	}

	fmt.Printf(
		"唱片总数：%d，平均价格：%.2f，最低价格：%.2f，最高价格：%.2f\n",
		stats.Count,
		stats.AveragePrice,
		stats.LowestPrice,
		stats.HighestPrice,
	)
	artistStatsList, err := getArtistStats(db)
	if err != nil {
		fmt.Println("艺术家统计失败：", err)
		return
	}

	fmt.Println("按艺术家统计：")

	for _, artistStats := range artistStatsList {
		fmt.Printf(
			"%s | %d 张 | 平均价格：%.2f 元\n",
			artistStats.Artist,
			artistStats.Count,
			artistStats.AveragePrice,
		)
	}
}
