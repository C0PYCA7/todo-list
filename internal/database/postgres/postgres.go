package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"todo-list/internal/config"
)

type Database struct {
	db *sql.DB
}

func New(dbConfig config.DatabaseConfig) (*Database, error) {
	const op = "database/postgres/postgres/New"

	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.Name,
	)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to open database: %w", op, err)
	}
	err = CreateTables(db)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create tables: %w", op, err)
	}

	return &Database{db: db}, err
}
