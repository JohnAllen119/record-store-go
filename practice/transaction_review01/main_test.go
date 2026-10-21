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
		t.Fatalf("open sql failed:%v", err)
	}
	t.Cleanup(func() {
		db.Close()
	})
	err = db.Ping()
	if err != nil {
		t.Fatalf("mysql connected failed:%v", err)
	}
	var dbName string
	err = db.QueryRow(`SELECT DATABASE()`).Scan(&dbName)
	if err != nil {
		t.Fatalf("sql query failed:%v", err)
	}
	if dbName != "record_store_tx_test" {
		t.Fatalf("dbName incorrect")
	}

	return db
}
func TestTransferStockRollback(t *testing.T) {
	db := openTestDB(t)
	ctx := context.Background()
	var initialStock int
	err := db.QueryRowContext(ctx, `
	SELECT stock
	FROM records
	WHERE record_id=1
	`).Scan(&initialStock)
	if err != nil {
		t.Fatalf("get initialstock failed:%v", err)
	}
	if initialStock <= 0 {
		t.Fatalf("initialStock<=0")
	}
	err = transferStock(ctx, db, 1)

	if err == nil {
		t.Fatalf("transferStock success:%v", err)
	}
	var mysqlErr *mysql.MySQLError
	if !errors.As(err, &mysqlErr) {
		t.Fatalf("mysql error:%v", err)
	}
	if mysqlErr.Number != 1146 {
		t.Fatalf("mysqlErr wrong:%v", err)
	}
	var finalStock int
	err = db.QueryRowContext(ctx, `
	SELECT stock
	FROM records
	WHERE record_id=1
	`).Scan(&finalStock)
	if err != nil {
		t.Fatalf("get finalStock failed:%v", err)
	}
	if finalStock != initialStock {
		t.Fatalf("finalStock != initialStock")
	}
}
func openSuccessTestDB(t *testing.T) *sql.DB {
	dsn := os.Getenv("RECORD_STORE_SUCCESS_TEST_DSN")
	if dsn == "" {
		t.Fatalf("get dsn failed")
	}
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("open mysql failed:%v", err)
	}
	t.Cleanup(func() {
		db.Close()
	})
	err = db.Ping()
	if err != nil {
		t.Fatalf("mysql connect failed:%v", err)
	}
	var dbName string
	err = db.QueryRow(`
	SELECT DATABASE()
	`).Scan(&dbName)
	if err != nil {
		t.Fatalf("db query failed:%v", err)
	}
	if dbName != "record_store_tx_success_test" {
		t.Fatalf("expect record_store_tx_success_test,get %s", dbName)
	}

	return db
}
func TestTransferStockCommit(t *testing.T) {
	db := openSuccessTestDB(t)
	ctx := context.Background()
	_, err := db.ExecContext(ctx, `
	INSERT INTO records(record_id,title,stock)
	VALUES (?,?,?)
	`, 999, "FXG", 5)
	if err != nil {
		t.Fatalf("insert record failed:%v", err)
	}
	t.Cleanup(func() {
		_, err := db.ExecContext(ctx, `
	DELETE FROM sales
	WHERE record_id=?
	`, 999)
		if err != nil {
			t.Errorf("delete sales failed: %v", err)
		}
		_, err = db.ExecContext(ctx, `
	DELETE FROM records
	WHERE record_id = ?
	`, 999)
		if err != nil {
			t.Errorf("delete records failed:%v", err)
		}

	})
	err = transferStock(ctx, db, 999)
	if err != nil {
		t.Fatalf("transferStock failed:%v", err)
	}
	var newStock int
	err = db.QueryRowContext(ctx, `
	SELECT stock
	FROM records
	WHERE record_id=?
	`, 999).Scan(&newStock)
	if err != nil {
		t.Fatalf("query mysql failed:%v", err)
	}
	if newStock != 4 {
		t.Fatalf("newStock!=4")
	}
	saleCount := 0
	totalQuantity := 0
	err = db.QueryRowContext(ctx, `
	SELECT COUNT(*), COALESCE(SUM(quantity), 0)
	FROM sales
	WHERE record_id = ?;
	`, 999).Scan(&saleCount, &totalQuantity)
	if err != nil {
		t.Fatalf("query sales failed: %v", err)
	}
	if saleCount != 1 || totalQuantity != 1 {
		t.Fatalf(
			"expected count=1, quantity=1, got count=%d, quantity=%d",
			saleCount, totalQuantity,
		)
	}
}

func TestTransferStockInsufficientStock(t *testing.T) {
	db := openSuccessTestDB(t)
	ctx := context.Background()
	_, err := db.ExecContext(ctx, `
	INSERT INTO records(record_id,title,stock)
	VALUES(?,?,?)
	`, 999, "FXG", 0)
	if err != nil {
		t.Fatalf("insert new record failed:%v", err)
	}
	t.Cleanup(func() {
		_, err = db.ExecContext(ctx, `
		DELETE FROM records
		WHERE record_id=?
		`, 999)
		if err != nil {
			t.Errorf("delete record failed:%v", err)
		}
	})
	err = transferStock(ctx, db, 999)
	if err == nil {
		t.Fatalf("transferStcok failed:%v", err)
	}
	stock := 0
	err = db.QueryRowContext(ctx, `
	SELECT stock
	FROM records
	WHERE record_id=?
	`, 999).Scan(&stock)
	if err != nil {
		t.Fatalf("query failed:%v", err)
	}
	if stock != 0 {
		t.Fatalf("stock != 0")
	}
	salesCount := 0
	stockQuantity := 0
	err = db.QueryRowContext(ctx, `
	SELECT COUNT(*),COALESCE(SUM(quantity),0)
	FROM sales
	WHERE record_id=?
	`, 999).Scan(&salesCount, &stockQuantity)
	if err != nil {
		t.Fatalf("query sales failed:%v", err)
	}
	if salesCount != 0 || stockQuantity != 0 {
		t.Fatalf("salesCount:%d, stockQuantity:%d",
			salesCount, stockQuantity)
	}

}
