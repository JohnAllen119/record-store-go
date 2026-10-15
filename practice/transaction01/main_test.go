package main

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
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
		t.Fatalf("mysql open failed")
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
