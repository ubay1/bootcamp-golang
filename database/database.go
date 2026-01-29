package database

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func InitDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	// PENTING untuk PgBouncer transaction mode
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(time.Minute * 5)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
