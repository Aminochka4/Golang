package model

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type Questionnaire struct {
	Id        string `json:"id"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	Topic     string `json:"topic"`
	Questions string `json:"questions"`
	UserId    string `json:"userId"`
}

type QuestionnaireModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
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
		return nil, err
	}
	return &questionnaire, nil
}

func (q QuestionnaireModel) Update(menu *Questionnaire) error {
	// Update a specific menu item in the database.
	query := `
		UPDATE questionnaire
		SET topic = $1, questions = $2,
		WHERE id = $4
		RETURNING updatedAt
		`
	args := []interface{}{menu.Topic, menu.Questions, menu.UserId, menu.Id}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return q.DB.QueryRowContext(ctx, query, args...).Scan(&menu.UpdatedAt)
}

func (q QuestionnaireModel) Delete(id int) error {
	// Delete a specific menu item from the database.
	query := `
		DELETE FROM questionnaire
		WHERE id = $1
		`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := q.DB.ExecContext(ctx, query, id)
	return err
}
