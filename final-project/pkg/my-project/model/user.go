package model

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Aminochka4/Golang/final-project/pkg/my-project/validator"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

var (
	ErrDuplicateEmail = errors.New("duplicate email")
)

var AnonymousUser = &User{}

type User struct {
	Id        int64     `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Name      string    `json:"name"`
	Surname   string    `json:"surname"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	Activated bool      `json:"activated"`
	Version   int       `json:"-"`
}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
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
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

func (u UserModel) Insert(user *User) error {
	query := `
			INSERT INTO users (name, surname, username, email, password, activated)
			VALUES($1, $2, $3, $4, $5, $6)
			RETURNING id, createdAt, version
			`
	args := []interface{}{user.Name, user.Surname, user.Username, user.Email, user.Password.hash, user.Activated}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	pqErr := `pq: duplicate key value violates unique constraint "users_email_key"`
	err := u.DB.QueryRowContext(ctx, query, args...).Scan(&user.Id, &user.CreatedAt, &user.Version)
	if err != nil {
		switch {
		case err.Error() == pqErr:
			return ErrDuplicateEmail
		default:
			return err
		}
	}

	return nil
}

func (u UserModel) GetAll() ([]*User, error) {
	query := `
		SELECT id, createdAt, name, surname, username, email, password, activated, version
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
		err := rows.Scan(&user.Id, &user.CreatedAt, &user.Name, &user.Surname,
			&user.Username, &user.Email, &user.Password.hash, &user.Activated, &user.Version)
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
		SELECT id, createdAt, name, surname, username, email, password, activated, version
		FROM users
		WHERE id = $1
		`
	var user User
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := u.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&user.Id, &user.CreatedAt,
		&user.Name, &user.Surname, &user.Username, &user.Email, &user.Password.hash, &user.Activated, &user.Version)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
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

	err := u.DB.QueryRowContext(ctx, query, args...).Scan(&user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
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

	query := `
		SELECT 
			users.id, users.createdAt, users.name, users.surname, users.username, 
			users.email, users.password, users.activated, users.version
		FROM	users
        INNER JOIN tokens
			ON users.id = tokens.user_id
        WHERE tokens.plaintext = $1 
			AND tokens.scope = $2
			AND tokens.expiry > $3
		`

	args := []interface{}{tokenPlaintext, tokenScope, time.Now()}

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := u.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.Id,
		&user.CreatedAt,
		&user.Name,
		&user.Surname,
		&user.Username,
		&user.Email,
		&user.Password.hash,
		&user.Activated,
		&user.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func ValidateUsername(v *validator.Validator, username string) {
	v.Check(username != "", "username", "must be provided")
	v.Check(validator.Matches(username, validator.UsernameRX), "username", "must be valid username format")
}

func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Name != "", "name", "must be provided")
	v.Check(len(user.Name) <= 500, "name", "must not be more than 500 bytes long")

	ValidateUsername(v, user.Username)

	if user.Password.plaintext != nil {
		ValidatePasswordPlaintext(v, *user.Password.plaintext)
	}

	if user.Password.hash == nil {
		// TODO: fix this panic
		panic("missing password hash for user")
	}
}

//tsis3

func (u *UserModel) GetByUsername(username string) (*User, error) {
	query := `
        SELECT id, createdAt, name, surname, username, email, password, activated, version
        FROM users
        WHERE username = $1
    `

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := u.DB.QueryRowContext(ctx, query, username).Scan(
		&user.Id,
		&user.CreatedAt,
		&user.Name,
		&user.Surname,
		&user.Username,
		&user.Email,
		&user.Password.hash,
		&user.Activated,
		&user.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (u *UserModel) GetBySurname() ([]*User, error) {
	query := `
        SELECT id, createdAt, name, surname, username, email, password, activated, version
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
		err := rows.Scan(&user.Id, &user.CreatedAt, &user.Name, &user.Surname,
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
        SELECT id, createdAt, name, surname, username, email, password, activated, version
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
		err := rows.Scan(&user.Id, &user.CreatedAt, &user.Name, &user.Surname,
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
