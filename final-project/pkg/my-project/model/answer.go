package model

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Aminochka4/Golang/final-project/pkg/my-project/validator"
	"log"
	"time"
)

type Answer struct {
	Id              string `json:"id"`
	CreatedAt       string `json:"createdAt"`
	UpdatedAt       string `json:"updatedAt"`
	QuestionnaireId string `json:"questionnaireId"`
	Answer          string `json:"answer"`
	UserId          int64  `json:"userId"`
}

type AnswerModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (a AnswerModel) GetAll() ([]*Answer, error) {
	query := `
		SELECT id, createdAt, updatedAt, questionnaireId, answer, userId
		FROM answer
		ORDER BY id
	`

	rows, err := a.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var answers []*Answer
	for rows.Next() {
		var answer Answer
		err := rows.Scan(&answer.Id, &answer.CreatedAt, &answer.UpdatedAt, &answer.QuestionnaireId, &answer.Answer, &answer.UserId)
		if err != nil {
			return nil, err
		}
		answers = append(answers, &answer)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return answers, nil
}

func (a AnswerModel) Insert(answer *Answer) error {
	// Insert a new menu item into the database.
	query := `
		INSERT INTO answer (questionnaireId, answer, userId) 
		VALUES ($1, $2, $3) 
		RETURNING id, createdAt, updatedAt
		`
	args := []interface{}{answer.QuestionnaireId, answer.Answer, answer.UserId}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return a.DB.QueryRowContext(ctx, query, args...).Scan(&answer.Id, &answer.CreatedAt, &answer.UpdatedAt)
}

func (a AnswerModel) Get(id int) (*Answer, error) {
	// Retrieve a specific menu item based on its ID.
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, createdAt, updatedAt, questionnaireId, answer, userId
		FROM answer
		WHERE id = $1
		`
	var answer Answer
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := a.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&answer.Id, &answer.CreatedAt, &answer.UpdatedAt, &answer.QuestionnaireId, &answer.Answer, &answer.UserId)
	if err != nil {
		return nil, fmt.Errorf("cannot retrive answer with id: %v, %w", id, err)
	}
	return &answer, nil
}

func (a AnswerModel) Update(answer *Answer) error {
	// Update a specific menu item in the database.
	query := `
		UPDATE answer
		SET questionnaireId = $1, answer = $2, userId = $3, updatedAt = CURRENT_TIMESTAMP
		WHERE id = $4 AND updatedAt = $5
		RETURNING updatedAt
		`
	args := []interface{}{answer.QuestionnaireId, answer.Answer, answer.UserId, answer.Id, answer.UpdatedAt}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return a.DB.QueryRowContext(ctx, query, args...).Scan(&answer.UpdatedAt)
}

func (a AnswerModel) Delete(id int) error {
	// Delete a specific menu item from the database.
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM answer
		WHERE id = $1
		`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := a.DB.ExecContext(ctx, query, id)
	return err
}

func (a AnswerModel) GetByQuestionnaire(questionnaireID int) ([]*Answer, error) {
	if questionnaireID < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
        SELECT id, createdAt, updatedAt, questionnaireId, answer, userId
        FROM answer
        WHERE questionnaireId = $1
    `

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := a.DB.QueryContext(ctx, query, questionnaireID)
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve answers for questionnaire with ID %d: %w", questionnaireID, err)
	}
	defer rows.Close()

	var answers []*Answer
	for rows.Next() {
		var answer Answer
		if err := rows.Scan(&answer.Id, &answer.CreatedAt, &answer.UpdatedAt, &answer.QuestionnaireId, &answer.Answer, &answer.UserId); err != nil {
			return nil, fmt.Errorf("cannot scan answer row: %w", err)
		}
		answers = append(answers, &answer)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading answer rows: %w", err)
	}

	return answers, nil
}

func ValidateAnswer(v *validator.Validator, answer *Answer) {
	// Check if the title field is empty.
	v.Check(answer.Answer != "", "topic", "must be provided")
	// Check if the description field is not more than 1000 characters.
	v.Check(len(answer.Answer) <= 1000, "questions", "must not be more than 1000 bytes long")
}
