package db

import "database/sql"

func InitDB() (*sql.DB, error) {
	conn := "user=biba dbname=postgres password=boba host=localhost sslmode=disable"
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
