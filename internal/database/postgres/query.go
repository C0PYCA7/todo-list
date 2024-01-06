package postgres

import (
	"database/sql"
	"fmt"
)

func CreateTables(db *sql.DB) error {

	const op = "database/postgres/query/CreateTables"

	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS username(
	    id SERIAL PRIMARY KEY,
	    name VARCHAR(255) NOT NULL
	);
	
	CREATE TABLE IF NOT EXISTS credential(
	    id SERIAL PRIMARY KEY,
	    id_user INT REFERENCES username(id), 
	    login VARCHAR(255),
	    password VARCHAR(255)
	);

	CREATE TABLE IF NOT EXISTS status(
	    id SERIAL PRIMARY KEY,
	    name varchar(255)
	);

	CREATE TABLE IF NOT EXISTS task(
	    id SERIAL PRIMARY KEY,
	    id_user INT REFERENCES username(id),
	    name VARCHAR(255),
	    id_status INT REFERENCES status(id),
	    description TEXT
	)
`)
	if err != nil {
		return fmt.Errorf("%s: failed to create query: %w", op, err)
	}

	return nil
}
