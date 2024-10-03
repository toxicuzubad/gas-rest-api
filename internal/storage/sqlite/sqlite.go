package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3" // Init sqllite3 driver.

	"gas-rest-api/internal/storage"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.NewStorage"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := db.Prepare(
		`CREATE TABLE IF NOT EXISTS guitar(
			id INTEGER PRIMARY KEY,
			manufacturer_name TEXT NOT NULL,
			model_name TEXT NOT NULL,
			description TEXT,
			serial_number TEXT); 
		CREATE INDEX IF NOT EXISTS idx_manufacturer_name ON guitar (manufacturer_name);
	`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveGuitar(manufacturerName string, modelName string, description string, serialNumber string) (int64, error) {
	const op = "storage.sqlite.SaveGuitar"

	stmt, err := s.db.Prepare("INSERT INTO guitar(manufacturer_name,model_name,description,serial_number) values(?,?,?,?)")
	if err != nil {
		return 0, fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	res, err := stmt.Exec(manufacturerName, modelName, description, serialNumber)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrGuitarExists)
		}

		return 0, fmt.Errorf("%s: execute statement: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
	}

	return id, nil
}
