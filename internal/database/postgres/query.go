package postgres

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"time"
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

	CREATE TABLE IF NOT EXISTS task(
	    id SERIAL PRIMARY KEY,
	    id_user INT REFERENCES username(id),
	    name VARCHAR(255),
	    description TEXT,
	    deadline timestamp
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

func (d *Database) SignIn(name, login, password string) (int64, error) {
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
		return 0, fmt.Errorf("%s: failed to insert signIn`s(%d) credential", op, userId)
	}
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("%s: failed to commit transaction: %w", op, err)
	}
	return userId, nil
}

func (d *Database) CreateTask(userId int64, name, description string, date time.Time) (int64, error) {
	const op = "database/postgres/query/CreateTask"

	var taskId int64
	//var count int

	//query := ``

	query := `INSERT INTO task(id_user, name, description, deadline)
				VALUES($1, $2, $3, $4) RETURNING task.id`

	err := d.db.QueryRow(query, userId, name, description, date).Scan(&taskId)
	if err != nil {
		return 0, fmt.Errorf("%s: failed to create task: %w", op, err)
	}
	return taskId, nil
}

type Task struct {
	Name        string
	Description string
	Deadline    time.Time
}

func (d *Database) ListOfTask(userId int64) ([]Task, error) {
	const op = "database/postgres/query/ListOfTask"

	query := "SELECT name, description, deadline " +
		"FROM task " +
		"WHERE id_user = $1"
	rows, err := d.db.Query(query, userId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var res []Task

	for rows.Next() {
		var list Task
		if err := rows.Scan(&list.Name, &list.Description, &list.Deadline); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		res = append(res, list)
	}
	return res, nil
}
