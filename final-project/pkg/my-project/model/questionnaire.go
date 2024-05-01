package model

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Aminochka4/Golang/final-project/pkg/my-project/validator"
	"log"
	"time"
)

type Questionnaire struct {
	Id        string `json:"id"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	Topic     string `json:"topic"`
	Questions string `json:"questions"`
	UserId    int64  `json:"userId"`
}

type QuestionnaireModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (q QuestionnaireModel) GetAll(topic string, filters Filters) ([]*Questionnaire, error) {
	// Формируем базовый запрос SQL
	query := `
		SELECT id, createdAt, updatedAt, topic, questions, userId
		FROM questionnaire
		WHERE ($1 = '' OR LOWER(topic) = LOWER($1))
	`

	// Добавляем сортировку в запрос, если указано значение Sort
	if filters.Sort != "" {
		query += " ORDER BY " + filters.sortColumn() + " " + filters.sortDirection()
	}

	// Добавляем параметры пагинации в запрос
	query += " LIMIT $2 OFFSET $3"

	rows, err := q.DB.Query(query, topic, filters.PageSize, (filters.Page-1)*filters.PageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questionnaires []*Questionnaire
	for rows.Next() {
		var questionnaire Questionnaire
		err := rows.Scan(&questionnaire.Id, &questionnaire.CreatedAt, &questionnaire.UpdatedAt, &questionnaire.Topic, &questionnaire.Questions, &questionnaire.UserId)
		if err != nil {
			return nil, err
		}
		questionnaires = append(questionnaires, &questionnaire)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return questionnaires, nil
}

func (q QuestionnaireModel) Insert(questionnaire *Questionnaire) error {
	// Insert a new menu item into the database.
	query := `
		INSERT INTO questionnaire (topic, questions, userId) 
		VALUES ($1, $2, $3) 
		RETURNING id, createdAt, updatedAt
		`
	args := []interface{}{questionnaire.Topic, questionnaire.Questions, questionnaire.UserId}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return q.DB.QueryRowContext(ctx, query, args...).Scan(&questionnaire.Id, &questionnaire.CreatedAt, &questionnaire.UpdatedAt)
}

func (q QuestionnaireModel) Get(id int) (*Questionnaire, error) {
	// Retrieve a specific menu item based on its ID.
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, createdAt, updatedAt, topic, questions, userId
		FROM questionnaire
		WHERE id = $1
		`
	var questionnaire Questionnaire
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := q.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&questionnaire.Id, &questionnaire.CreatedAt, &questionnaire.UpdatedAt, &questionnaire.Topic, &questionnaire.Questions, &questionnaire.UserId)
	if err != nil {
		return nil, fmt.Errorf("cannot retrive questionnaire with id: %v, %w", id, err)
	}
	return &questionnaire, nil
}

func (q QuestionnaireModel) Update(questionnaire *Questionnaire) error {
	// Update a specific menu item in the database.
	query := `
		UPDATE questionnaire
		SET topic = $1, questions = $2, userId = $3, updatedAt = CURRENT_TIMESTAMP
		WHERE id = $4 AND updatedAt = $5
		RETURNING updatedAt
		`
	args := []interface{}{questionnaire.Topic, questionnaire.Questions, questionnaire.UserId, questionnaire.Id, questionnaire.UpdatedAt}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return q.DB.QueryRowContext(ctx, query, args...).Scan(&questionnaire.UpdatedAt)
}

func (q QuestionnaireModel) Delete(id int) error {
	// Delete a specific menu item from the database.
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM questionnaire
		WHERE id = $1
		`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := q.DB.ExecContext(ctx, query, id)
	return err
}

func ValidateQuestionnaire(v *validator.Validator, questionnaire *Questionnaire) {
	// Check if the title field is empty.
	v.Check(questionnaire.Topic != "", "topic", "must be provided")
	// Check if the title field is not more than 100 characters.
	v.Check(len(questionnaire.Topic) <= 100, "topic", "must not be more than 100 bytes long")
	// Check if the description field is not more than 1000 characters.
	v.Check(len(questionnaire.Questions) <= 1000, "questions", "must not be more than 1000 bytes long")
}
