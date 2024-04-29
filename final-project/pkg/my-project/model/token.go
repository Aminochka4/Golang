package model

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base32"
	"errors"
	"github.com/Aminochka4/Golang/final-project/pkg/my-project/validator"
	"log"
	"time"
)

const (
	ScopeActivation     = "activation"
	ScopeAuthentication = "authentication"
)

type (
	Token struct {
		Plaintext string    `json:"token"`
		UserID    int64     `json:"-"`
		Expiry    time.Time `json:"expiry"`
		Scope     string    `json:"-"`
	}

	TokenModel struct {
		DB       *sql.DB
		InfoLog  *log.Logger
		ErrorLog *log.Logger
	}
)

func (m TokenModel) Parse(tokenString string) (*Token, error) {
	// Напишите SQL-запрос для поиска токена по его строковому представлению
	query := `
		SELECT plaintext, user_id, expiry, scope
		FROM tokens
		WHERE plaintext = $1
	`

	// Выполните SQL-запрос и получите результат
	var token Token
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, tokenString).Scan(&token.Plaintext, &token.UserID, &token.Expiry, &token.Scope)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	// Установите строковое представление токена
	if token.Plaintext == tokenString {
		return &token, nil
	}

	return nil, ErrRecordNotFound
}

func (m TokenModel) New(userID int64, ttl time.Duration, scope string) (*Token, error) {
	token, err := generateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}

	err = m.Insert(token)
	return token, err

}

func (m TokenModel) Insert(token *Token) error {
	query := `
		INSERT INTO tokens (plaintext, user_id, expiry, scope)
		VALUES ($1, $2, $3, $4)
		`

	args := []interface{}{token.Plaintext, token.UserID, token.Expiry, token.Scope}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, args...)
	return err
}

func (m TokenModel) DeleteAllForUser(scope string, userID int64) error {
	query := `
		DELETE FROM tokens
		WHERE scope = $1 AND user_id = $2
		`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, scope, userID)
	return err
}

func generateToken(userID int64, ttl time.Duration, scope string) (*Token, error) {
	token := &Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}

	randomBytes := make([]byte, 16)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	return token, nil
}

func ValidateTokenPlaintext(v *validator.Validator, tokenPlaintext string) {
	v.Check(tokenPlaintext != "", "token", "must be provided")
	v.Check(len(tokenPlaintext) == 26, "token", "must be 26 bytes long")
}
