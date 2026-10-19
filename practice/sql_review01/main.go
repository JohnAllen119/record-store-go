package main
import(
	"database/sql"
	"context"
)

type Record struct{
	ID int
	Title string
	Stock int
}

func findRecordByID(
    ctx context.Context,
    db *sql.DB,
    id int,
) (Record, error) {
	var record Record
	err := db.QueryRowContext(ctx,`
	SELECT record_id,title,stock
	FROM records
	WHERE record_id=?
	`,id).Scan(&record.ID,&record.Title,&record.Stock)
	if err != nil{
		return Record{},err
	}
	return record,nil
}
func updateRecordStock(
	ctx context.Context,
	db *sql.DB,
	id int,
	newStock int,
) (bool,error){
	result,err := db.ExecContext(ctx,`
	UPDATE records
	SET stock=?
	WHERE record_id=?
	`,newStock,id)
	if err != nil{
		return false,err
	}
	affected,err:=result.RowsAffected()
	if err!=nil{
		return false,err
	}
	if affected == 0{
		return false,nil
	}
	return true,nil
}
