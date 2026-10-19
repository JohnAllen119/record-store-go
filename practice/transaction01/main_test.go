package main

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"os"
	"testing"
)

func openTestDB(t *testing.T) *sql.DB {
	dsn := os.Getenv("RECORD_STORE_TEST_DSN")
	if dsn == "" {
		t.Fatalf("get dsn failed")
	}
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("mysql open failed: %v", err)
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		t.Fatalf("mysql connect failed: %v", err)
	}
	var dbName string
	err = db.QueryRow("SELECT DATABASE()").Scan(&dbName)
	if err != nil {
		db.Close()
		t.Fatalf("check database failed:%v", err)
	}
	if dbName != "record_store_tx_test" {
		db.Close()
		t.Fatalf("wrong test database: %s", dbName)
	}
	return db
}

func TestDatabaseConnection(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

}
func TestSellRecordRollback(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	ctx := context.Background()
	var id int
	var stock int
	err := db.QueryRowContext(ctx, `
		SELECT record_id,stock
		FROM records
		WHERE record_id=1
		`).Scan(&id, &stock)
	if err != nil {
		t.Fatalf("query initial stock failed: %v", err)
	}
	if id != 1 || stock <= 0 {
		t.Fatalf("invalid initial record: id=%d, stock=%d", id, stock)
	}
	var tablecount int
	err = db.QueryRowContext(ctx, `
	SELECT COUNT(*)
FROM information_schema.tables
WHERE table_schema = DATABASE()
  AND table_name = 'sales';
  `).Scan(&tablecount)
	if err != nil {
		t.Fatalf("check sales table failed: %v", err)
	}
	if tablecount > 0 {
		t.Skip("sales table exists")
	}
	err = sellRecord(ctx, db, 1)
	if err == nil {
		t.Fatalf("sales shouldn't exist")
	}
	var mysqlErr *mysql.MySQLError
	if !errors.As(err, &mysqlErr) {
		t.Fatalf("expected MySQL error, got %v", err)
	}
	if mysqlErr.Number != 1146 {
		t.Fatalf("expected 1146,got %d", mysqlErr.Number)
	}
	var newstock int
	err = db.QueryRowContext(ctx, `
		SELECT record_id,stock
		FROM records
		WHERE record_id=1
		`).Scan(&id, &newstock)
	if err != nil {
		t.Fatalf("query final stock failed: %v", err)
	}
	if newstock != stock {
		t.Fatalf("expected stock %d, got %d", stock, newstock)
	}
}
func openSuccessTestDB(t *testing.T) *sql.DB {
	dsn := os.Getenv("RECORD_STORE_SUCCESS_TEST_DSN")
	if dsn == "" {
		t.Fatalf("get dsn failed")
	}
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("mysql open failed: %v", err)
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		t.Fatalf("mysql connect failed: %v", err)
	}
	var dbName string
	err = db.QueryRow(`SELECT DATABASE()`).Scan(&dbName)
	if err != nil {
		defer db.Close()
		t.Fatalf("select database failed: %s", dbName)
	}
	if dbName != "record_store_tx_success_test" {
		defer db.Close()
		t.Fatalf("check mysql failed: %v", err)
	}
	return db
}
func TestSuccessDatabaseConnection(t *testing.T) {
	db := openSuccessTestDB(t)
	defer db.Close()
}
func TestSellRecordCommit(t *testing.T) {
	db := openSuccessTestDB(t)
	defer db.Close()
	ctx := context.Background()
	id := 900001
	var Stock int
	_, err := db.ExecContext(ctx, `
	INSERT INTO records(record_id,title,stock)
	VALUES (?,?,?)
	`, id, "Test Record", 5)
	if err != nil {
		t.Fatalf("insert failed %v", err)
	}
	defer func() {
		if _, err := db.ExecContext(ctx,
			`DELETE FROM sales WHERE record_id=?
		`, id); err != nil {
			t.Errorf("clean sales failed:%v", err)
		}
		if _, err := db.ExecContext(ctx, `
		DELETE FROM records WHERE record_id=?
		`, id); err != nil {
			t.Errorf("cleanup record failed: %v", err)
		}
	}()
	err = db.QueryRowContext(ctx, `
	SELECT stock
	FROM records
	WHERE record_id = ?
	`, id).Scan(&Stock)
	if err != nil {
		t.Fatalf("select failed:%v", err)
	}

	err = sellRecord(ctx, db, id)
	if err != nil {
		t.Fatalf("insert failed :%v", err)
	}
	var newStock int
	err = db.QueryRowContext(ctx, `
	SELECT stock
	FROM records
	WHERE record_id = ? `, id).Scan(&newStock)
	if err != nil {
		t.Fatalf("select failed:%v", err)
	}
	if newStock+1 != Stock {
		t.Fatalf("expect stock%d,got %d", Stock-1, newStock)
	}
	saleCount := 0
	totalQuantity := 0
	err = db.QueryRowContext(ctx, `
	SELECT COUNT(*),COALESCE(SUM(quantity),0)
	FROM sales
	WHERE record_id=?
	`, id).Scan(&saleCount, &totalQuantity)
	if err != nil {
		t.Fatalf("select failed %v", err)
	}
	if saleCount != 1 || totalQuantity != 1 {
		t.Fatalf("expected salecount:%d,totalQuantity:%v", saleCount, totalQuantity)
	}

}
