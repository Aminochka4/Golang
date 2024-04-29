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

func (q QuestionnaireModel) GetAll(title string, from, to int, filters Filters) ([]*Questionnaire, Metadata, error) {

	// Retrieve all menu items from the database.
	query := fmt.Sprintf(
		`
		SELECT count(*) OVER(), id, createdAt, updatedAt, topic, questions, userId
		FROM questionnaire
		WHERE (LOWER(topic) = LOWER($1) OR $1 = '')
		ORDER BY %s %s, id ASC
		LIMIT $4 OFFSET $5
		`,
		filters.sortColumn(), filters.sortDirection())

	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Organize our four placeholder parameter values in a slice.
	args := []interface{}{title, from, to, filters.limit(), filters.offset()}

	// log.Println(query, title, from, to, filters.limit(), filters.offset())
	// Use QueryContext to execute the query. This returns a sql.Rows result set containing
	// the result.
	rows, err := q.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	// Importantly, defer a call to rows.Close() to ensure that the result set is closed
	// before GetAll returns.
	defer func() {
		if err := rows.Close(); err != nil {
			q.ErrorLog.Println(err)
		}
	}()

	// Declare a totalRecords variable
	totalRecords := 0

	var questionnaires []*Questionnaire
	for rows.Next() {
		var questionnaire Questionnaire
		err := rows.Scan(&totalRecords, &questionnaire.Id, &questionnaire.CreatedAt, &questionnaire.UpdatedAt, &questionnaire.Topic, &questionnaire.Questions, &questionnaire.UserId)
		if err != nil {
			return nil, Metadata{}, err
		}

		// Add the Movie struct to the slice
		questionnaires = append(questionnaires, &questionnaire)
	}

	// When the rows.Next() loop has finished, call rows.Err() to retrieve any error
	// that was encountered during the iteration.
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	// Generate a Metadata struct, passing in the total record count and pagination parameters
	// from the client.
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	// If everything went OK, then return the slice of the movies and metadata.
	return questionnaires, metadata, nil
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

func ValidateMenu(v *validator.Validator, questionnaire *Questionnaire) {
	// Check if the title field is empty.
	v.Check(questionnaire.Topic != "", "topic", "must be provided")
	// Check if the title field is not more than 100 characters.
	v.Check(len(questionnaire.Topic) <= 100, "topic", "must not be more than 100 bytes long")
	// Check if the description field is not more than 1000 characters.
	v.Check(len(questionnaire.Questions) <= 1000, "questions", "must not be more than 1000 bytes long")
}
