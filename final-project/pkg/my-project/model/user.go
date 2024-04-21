package model

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"github.com/Aminochka4/Golang/final-project/pkg/my-project/validator"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

type User struct {
	Id        int64     `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Name      string    `json:"name"`
	Surname   string    `json:"surname"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	Activated bool      `json:"activated"`
	Version   int       `json:"-"`
}
type UserModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

type password struct {
	plaintext *string
	hash      []byte
}

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}
	p.plaintext = &plaintextPassword
	p.hash = hash
	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		return false, err
	}

	return true, nil
}

func (u UserModel) Insert(user *User) error {
	query := `
			INSERT INTO users (name, surname, username, email, password, activated)
			VALUES($1, $2, $3, $4, $5, $6)
			RETURNING id, createdAt, updatedAt, version
			`
	args := []interface{}{user.Name, user.Surname, user.Username, user.Email, user.Password.hash, user.Activated}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return u.DB.QueryRowContext(ctx, query, args...).Scan(&user.Id, &user.CreatedAt, &user.UpdatedAt, &user.Version)
}

func (u UserModel) GetAll() ([]*User, error) {
	query := `
		SELECT id, createdAt, updatedAt, name, surname, username, email, password, activated, version
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
			&user.Name, &user.Surname, &user.Username, &user.Email, &user.Password.hash, &user.Activated, &user.Version)
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
		SELECT id, createdAt, updatedAt, name, surname, username, email, password, activated, version
		FROM users
		WHERE id = $1
		`
	var user User
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := u.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&user.Id, &user.CreatedAt, &user.UpdatedAt,
		&user.Name, &user.Surname, &user.Username, &user.Email, &user.Password.hash, &user.Activated, &user.Version)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u UserModel) Update(user *User) error {
	query := `
		UPDATE users
		SET  name = $1, surname = $2, username = $3, email = $4, password = $5, activated = $6, version = version + 1
		WHERE id = $7 AND version = $8
		RETURNING version
		`

	args := []interface{}{user.Name, user.Surname, user.Username, user.Email, user.Password.hash, user.Activated, user.Id, user.Version}
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

func (u UserModel) GetForToken(tokenScope, tokenPlaintext string) (*User, error) {
	tokenHash := sha256.Sum256([]byte(tokenPlaintext))

	query := `
		SELECT 
			user.id, user.createdAt, user.updatedAt, user.name, user.surname, user.username, 
			user.email, user.password, user.activated, user.version
		FROM	users
        INNER JOIN tokens
			ON users.id = tokens.user_id
        WHERE tokens.hash = $1 
			AND tokens.scope = $2
			AND tokens.expiry > $3
		`

	args := []interface{}{tokenHash[:], tokenScope, time.Now()}

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := u.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.Id,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Name,
		&user.Surname,
		&user.Username,
		&user.Email,
		&user.Password.hash,
		&user.Activated,
		&user.Version,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be valid email address")
}

func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Name != "", "name", "must be provided")
	v.Check(len(user.Name) <= 500, "name", "must not be more than 500 bytes long")

	ValidateEmail(v, user.Email)

	if user.Password.plaintext != nil {
		ValidatePasswordPlaintext(v, *user.Password.plaintext)
	}

	if user.Password.hash == nil {
		// TODO: fix this panic
		panic("missing password hash for user")
	}
}

//tsis3

func (u *UserModel) GetByName(name string) ([]*User, error) {
	query := `
        SELECT id, createdAt, updatedAt, name, surname, username, email, password, activated, version
        FROM users
        WHERE name = $1
    `
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := u.DB.QueryContext(ctx, query, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user := &User{}
		err := rows.Scan(&user.Id, &user.CreatedAt, &user.UpdatedAt, &user.Name, &user.Surname,
			&user.Username, &user.Email, &user.Password, &user.Activated, &user.Version)
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

func (u *UserModel) GetBySurname() ([]*User, error) {
	query := `
        SELECT id, createdAt, updatedAt, name, surname, username, email, password, activated, version
        FROM users
        ORDER BY surname
    `
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := u.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user := &User{}
		err := rows.Scan(&user.Id, &user.CreatedAt, &user.UpdatedAt, &user.Name, &user.Surname,
			&user.Username, &user.Email, &user.Password, &user.Activated, &user.Version)
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

func (u *UserModel) GetUsersWithPagination(limit, offset int) ([]*User, error) {
	query := `
        SELECT id, createdAt, updatedAt, name, surname, username, email, password, activated, version
        FROM users
        ORDER BY id
        LIMIT $1 OFFSET $2
    `
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := u.DB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user := &User{}
		err := rows.Scan(&user.Id, &user.CreatedAt, &user.UpdatedAt, &user.Name, &user.Surname,
			&user.Username, &user.Email, &user.Password, &user.Activated, &user.Version)
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
