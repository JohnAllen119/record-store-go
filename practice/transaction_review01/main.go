package main

import (
	"context"
	"database/sql"
	"fmt"
)

func transferStock(ctx context.Context, db *sql.DB,recordID int) error {
	Tx, err := db.BeginTx(ctx, nil)

	if err != nil {
		return fmt.Errorf("tx begin failed:%w", err)
	}
	defer Tx.Rollback()
	result, err := Tx.ExecContext(ctx, `
	UPDATE records
	SET stock=stock-1
	WHERE record_id=? AND stock>0
	`,recordID)
	if err != nil {
		return fmt.Errorf("update records stock  failed:%w", err)
	}
	ra, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get affected rows failed:%w", err)
	}
	if ra == 0 {
		return fmt.Errorf("no matched record to update")
	}
	result, err = Tx.ExecContext(ctx, `
	INSERT INTO sales(record_id,quantity)
	VALUES (?,1)
	`,recordID)
	if err != nil {
		return fmt.Errorf("insert salerecord failed:%w", err)
	}
	ra, err = result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get affected rows failed:%w", err)
	}
	if ra == 0 {
		return fmt.Errorf("no matched record to insert")
	}
	err = Tx.Commit()
	if err != nil {
		return fmt.Errorf("commit failed:%w", err)
	}
	return nil
}
