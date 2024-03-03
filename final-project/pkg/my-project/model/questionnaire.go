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
	UserId    uint   `json:"user_id"`
}

type QuestionnaireModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (m QuestionnaireModel) Insert(questionnaire *Questionnaire) error {
	// Insert a new menu item into the database.
	query := `
		INSERT INTO menu (topic, questions, user_id) 
		VALUES ($1, $2, $3) 
		RETURNING id, created_at, updated_at
		`
	args := []interface{}{questionnaire.Topic, questionnaire.Questions, questionnaire.UserId}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&questionnaire.Id, &questionnaire.CreatedAt, &questionnaire.UpdatedAt)
}

func (m QuestionnaireModel) Get(id int) (*Questionnaire, error) {
	// Retrieve a specific menu item based on its ID.
	query := `
		SELECT id, created_at, updated_at, title, description, nutrition_value
		FROM menu
		WHERE id = $1
		`
	var menu Questionnaire
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&menu.Id, &menu.CreatedAt, &menu.UpdatedAt, &menu.Topic, &menu.Questions, &menu.UserId)
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

func (m QuestionnaireModel) Update(menu *Questionnaire) error {
	// Update a specific menu item in the database.
	query := `
		UPDATE menu
		SET title = $1, description = $2, nutrition_value = $3
		WHERE id = $4
		RETURNING updated_at
		`
	args := []interface{}{menu.Topic, menu.Questions, menu.UserId, menu.Id}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&menu.UpdatedAt)
}

func (m QuestionnaireModel) Delete(id int) error {
	// Delete a specific menu item from the database.
	query := `
		DELETE FROM menu
		WHERE id = $1
		`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}
