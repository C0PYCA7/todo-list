package postgres

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
)

func CreateTables(db *sql.DB) error {

	const op = "database/postgres/query/CreateTables"

	query := `
	CREATE TABLE IF NOT EXISTS username(
	    id SERIAL PRIMARY KEY,
	    name VARCHAR(255) NOT NULL UNIQUE
	);
	
	CREATE TABLE IF NOT EXISTS credential(
	    id SERIAL PRIMARY KEY,
	    id_user INT REFERENCES username(id), 
	    login VARCHAR(255) UNIQUE ,
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
	);
`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("%s: failed to create query: %w", op, err)
	}

	return nil
}

func (d *Database) CheckUser(login, password string) (int64, error) {

	const op = "database/postgres/query/CheckUser"

	var userId int64

	query := `SELECT id_user FROM credential WHERE login = $1 AND password = $2`

	err := d.db.QueryRow(query, login, password).Scan(&userId)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, ErrUserNotFound)
	}

	return userId, nil
}

func (d *Database) CreateUser(name, login, password string) (int64, error) {
	const op = "database/postgres/query/CreateUser"

	var userId int64

	tx, err := d.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to begin transacrion: %w", op, err)
	}
	defer tx.Rollback()

	queryUsername := `INSERT INTO username (name) VALUES($1) RETURNING id`

	err = tx.QueryRow(queryUsername, name).Scan(&userId)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				return 0, fmt.Errorf("%s: %w", op, ErrUserNameExists)
			}
		}
		return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
	}

	queryCredential := "INSERT INTO credential(id_user, login, password) VALUES($1, $2, $3)"

	_, err = tx.Exec(queryCredential, userId, login, password)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				return 0, fmt.Errorf("%s: %w", op, ErrUserLoginExists)
			}
		}
		return 0, fmt.Errorf("%s: failed to insert user`s(%d) credential", op, userId)
	}
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("%s: failed to commit transaction: %w", op, err)
	}
	return userId, nil
}
