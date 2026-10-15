package main

import (
	"context"
	"database/sql"
	"fmt"
)

func sellRecord(ctx context.Context, db *sql.DB, recordID int) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	result, err := tx.ExecContext(ctx, `UPDATE records SET stock = stock-1
	 WHERE record_id =? AND stock >0 `, recordID)
	if err != nil {
		return err
	}
	affectedRows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get affected rows failed: %w", err)
	}
	if affectedRows == 0 {
		return fmt.Errorf("record not found or insufficient stock")
	}
	_, err = tx.ExecContext(ctx, "INSERT INTO sales (record_id,quantity) VALUES(?,1)", recordID)
	if err != nil {
		return err
	}
	return tx.Commit()
}
