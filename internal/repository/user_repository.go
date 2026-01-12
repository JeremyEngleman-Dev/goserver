package repository

import (
	"database/sql"
	"errors"

	"example/goserver/internal/model"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateTable() error {
	_, err := r.db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			age INTEGER NOT NULL
		);
	`)
	return err
}

func (r *Repository) ListEmployees() ([]model.Employee, error) {
	rows, err := r.db.Query(`SELECT id, name, age FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.Employee
	for rows.Next() {
		var u model.Employee
		if err := rows.Scan(&u.ID, &u.Name, &u.Age); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

var ErrNotFound = errors.New("record not found")

func (r *Repository) GetEmployeeByID(id int64) (*model.Employee, error) {
	var u model.Employee

	err := r.db.QueryRow(
		`SELECT id, name, age FROM users WHERE id = ?`,
		id,
	).Scan(&u.ID, &u.Name, &u.Age)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &u, nil
}

func (r *Repository) CreateEmployee(u *model.Employee) error {
	result, err := r.db.Exec(
		`INSERT INTO users (name, age) VALUES (?, ?)`,
		u.Name,
		u.Age,
	)
	if err != nil {
		return err
	}

	u.ID, _ = result.LastInsertId()
	return nil
}
