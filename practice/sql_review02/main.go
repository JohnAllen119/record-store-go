package main

import (
	"context"
	"database/sql"
)

type Record struct {
	ID     int
	Title  string
	Artist string
	Stock  int
}

func findRecordByID(ctx context.Context, db *sql.DB, id int) (Record, error) {
	var record Record
	err := db.QueryRowContext(ctx, `
	SELECT record_id,title,artist,stock
	FROM records
	WHERE record_id=?
	`, id).Scan(&record.ID, &record.Title, &record.Artist, &record.Stock)
	if err == sql.ErrNoRows {
		return Record{}, err
	}
	if err != nil {
		return Record{}, err
	}
	return record, nil
}
func updateRecordStock(ctx context.Context, db *sql.DB, id int, newStock int) (int64, error) {
	result, err := db.ExecContext(ctx, `
	UPDATE records
	SET stock=?
	WHERE record_id=?
	`, newStock, id)
	if err != nil {
		return 0, err
	}
	cnt, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return cnt, nil

}
