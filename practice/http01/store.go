package main

import "database/sql"

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
	records := make([]Record, 0)
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
func getRecordsByArtist(db *sql.DB, artist string) ([]Record, error) {
	rows, err := db.Query(`
	SELECT record_id,title,artist,price
	FROM records
	WHERE artist=?
	ORDER BY record_id
	`, artist)
	if err != nil {
		return []Record{}, err
	}
	defer rows.Close()
	records := make([]Record, 0)
	for rows.Next() {
		var record Record
		err = rows.Scan(&record.ID, &record.Title, &record.Artist, &record.Price)
		if err != nil {
			return []Record{}, err
		}
		records = append(records, record)
	}
	if err = rows.Err(); err != nil {
		return []Record{}, err
	}
	return records, nil
}
func getRecordsByMinPrice(db *sql.DB, minPrice int) ([]Record, error) {
	records := make([]Record, 0)
	rows, err := db.Query(`
	SELECT record_id,title,artist,price
	FROM records
	WHERE price >=?
	ORDER BY record_id
	`, minPrice)
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
	if err = rows.Err(); err != nil {
		return []Record{}, err
	}
	return records, nil
}

func getRecordsByArtistAndMinPrice(
	db *sql.DB, artist string, minPrice int,
) ([]Record, error) {
	records := make([]Record, 0)
	rows, err := db.Query(`
	SELECT record_id,title,artist,price
	FROM records
	WHERE artist = ? AND price >= ?
	ORDER BY record_id
	`, artist, minPrice)
	if err != nil {
		return []Record{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var record Record
		err = rows.Scan(&record.ID, &record.Title, &record.Artist, &record.Price)
		if err != nil {
			return []Record{}, err
		}
		records = append(records, record)
	}
	if err = rows.Err(); err != nil {
		return []Record{}, err
	}
	return records, nil
}
