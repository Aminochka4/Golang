package model

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type User struct {
	Id        string `json:"id"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	Name      string `json:"name"`
	Surname   string `json:"surname"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}
type UserModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (m *UserModel) GetByName(name string) ([]*User, error) {
	query := `
        SELECT id, createdAt, updatedAt, name, surname, username, email, password
        FROM users
        WHERE name = $1
    `
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user := &User{}
		err := rows.Scan(&user.Id, &user.CreatedAt, &user.UpdatedAt,
			&user.Name, &user.Surname, &user.Username, &user.Email, &user.Password)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (m *UserModel) GetBySurname() ([]*User, error) {
	query := `
        SELECT id, createdAt, updatedAt, name, surname, username, email, password
        FROM users
        ORDER BY surname
    `
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user := &User{}
		err := rows.Scan(&user.Id, &user.CreatedAt, &user.UpdatedAt,
			&user.Name, &user.Surname, &user.Username, &user.Email, &user.Password)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// GetUsersWithPagination извлекает пользователей из базы данных с учетом лимита и смещения.
func (m *UserModel) GetUsersWithPagination(limit, offset int) ([]*User, error) {
	query := `
        SELECT id, createdAt, updatedAt, name, surname, username, email, password
        FROM users
        ORDER BY id
        LIMIT $1 OFFSET $2
    `
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user := &User{}
		err := rows.Scan(&user.Id, &user.CreatedAt, &user.UpdatedAt,
			&user.Name, &user.Surname, &user.Username, &user.Email, &user.Password)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (u UserModel) Insert(user *User) error {
	query := `
			INSERT INTO users (name, surname, username, email, password)
			VALUES($1, $2, $3, $4, $5)
			RETURNING id, createdAt, updatedAt
			`
	args := []interface{}{user.Name, user.Surname, user.Username, user.Email, user.Password}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	return u.DB.QueryRowContext(ctx, query, args...).Scan(&user.Id, &user.CreatedAt, &user.UpdatedAt)
}

func (u UserModel) GetAll() ([]*User, error) {
	query := `
		SELECT id, createdAt, updatedAt, name, surname, username, email, password
		FROM users
		ORDER BY id
	`

	rows, err := u.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.CreatedAt, &user.UpdatedAt,
			&user.Name, &user.Surname, &user.Username, &user.Email, &user.Password)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (u UserModel) GetById(id int) (*User, error) {
	query := `
		SELECT id, createdAt, updatedAt, name, surname, username, email, password
		FROM users
		WHERE id = $1
		`
	var user User
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := u.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&user.Id, &user.CreatedAt, &user.UpdatedAt,
		&user.Name, &user.Surname, &user.Username, &user.Email, &user.Password)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u UserModel) Update(user *User) error {
	query := `
		UPDATE users
		SET  name = $1, surname = $2, username = $3, email = $4, password = $5
		WHERE id = $6
		RETURNING updatedAt
		`

	args := []interface{}{user.Name, user.Surname, user.Username, user.Email, user.Password, user.Id}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return u.DB.QueryRowContext(ctx, query, args...).Scan(&user.UpdatedAt)
}

func (u UserModel) Delete(id int) error {
	query := `
		DELETE FROM users
		WHERE id = $1
		`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	_, err := u.DB.ExecContext(ctx, query, id)

	return err
}
