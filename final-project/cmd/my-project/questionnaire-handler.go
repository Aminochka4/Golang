package main

import (
	"errors"
	"github.com/Aminochka4/Golang/final-project/pkg/my-project/model"
	"github.com/Aminochka4/Golang/final-project/pkg/my-project/validator"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func (app *application) respondWithError(w http.ResponseWriter, code int, message string) {
	app.respondWithJson(w, code, map[string]string{"error": message})
}

func (app *application) getAllQuestionnairesHandler(w http.ResponseWriter, r *http.Request) {
	v := validator.New()
	// Извлекаем значение параметра topic из URL
	topic := r.URL.Query().Get("topic")

	// Извлекаем значение параметра сортировки (Sort) из URL
	sort := r.URL.Query().Get("sort")

	// Извлекаем параметры пагинации из URL
	page := app.readInt(r.URL.Query(), "page", 1, v)
	pageSize := app.readInt(r.URL.Query(), "page_size", 10, v)

	// Создаем экземпляр структуры Filters и устанавливаем параметры сортировки и пагинации
	filters := model.Filters{
		Sort:         sort,
		SortSafeList: []string{"id", "createdAt", "updatedAt", "topic", "userId"}, // Перечислите допустимые поля для сортировки
		Page:         page,
		PageSize:     pageSize,
	}

	// Вызываем функцию GetAll с переданными значениями topic и filters
	questionnaires, err := app.models.Questionnaires.GetAll(topic, filters)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "Failed to fetch questionnaires")
		return
	}

	app.respondWithJson(w, http.StatusOK, questionnaires)
}

func (app *application) extractToken(r *http.Request) (string, error) {
	// Извлечь токен из заголовка Authorization
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header is missing")
	}
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("invalid authorization header format")
	}

	// Вернуть токен
	return parts[1], nil
}

func (app *application) getUserIdFromToken(tokenString string) (int64, error) {
	// Распарсить токен и получить userId
	token, err := app.models.Tokens.Parse(tokenString)
	if err != nil {
		return 0, err
	}
	return token.UserID, nil
}

func (app *application) createQuestionnaireHandler(w http.ResponseWriter, r *http.Request) {
	tokenString, err := app.extractToken(r)
	if err != nil {
		app.respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Получить userId из токена
	userID, err := app.getUserIdFromToken(tokenString)
	//fmt.Println("createQuestionnaireHandler called")
	if err != nil {
		app.respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var input struct {
		Topic     string `json:"topic"`
		Questions string `json:"questions"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		log.Println(err)
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}

	questionnaire := &model.Questionnaire{
		Topic:     input.Topic,
		Questions: input.Questions,
		UserId:    userID,
	}

	err = app.models.Questionnaires.Insert(questionnaire)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJson(w, http.StatusCreated, questionnaire)
}

func (app *application) getQuestionnaireHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["questionnaireId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid questionnaire ID")
		return
	}

	questionnaire, err := app.models.Questionnaires.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	app.respondWithJson(w, http.StatusOK, questionnaire)
}

func (app *application) updateQuestionnaireHandler(w http.ResponseWriter, r *http.Request) {
	tokenString, err := app.extractToken(r)
	if err != nil {
		app.respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Получить userId из токена
	userID, err := app.getUserIdFromToken(tokenString)
	//fmt.Println("createQuestionnaireHandler called")
	if err != nil {
		app.respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	param := vars["questionnaireId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid questionnaire ID")
		return
	}

	questionnaire, err := app.models.Questionnaires.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	var input struct {
		Topic     *string `json:"topic"`
		Questions *string `json:"questions"`
	}

	if userID != questionnaire.UserId {
		app.respondWithError(w, http.StatusInternalServerError, "Cannot update other's questionnaire")
		return
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if input.Topic != nil {
		questionnaire.Topic = *input.Topic
	}

	if input.Questions != nil {
		questionnaire.Questions = *input.Questions
	}

	err = app.models.Questionnaires.Update(questionnaire)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
	}

	app.respondWithJson(w, http.StatusOK, questionnaire)
}

func (app *application) deleteQuestionnaireHandler(w http.ResponseWriter, r *http.Request) {
	tokenString, err := app.extractToken(r)
	if err != nil {
		app.respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Получить userId из токена
	userID, err := app.getUserIdFromToken(tokenString)
	//fmt.Println("createQuestionnaireHandler called")
	if err != nil {
		app.respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	vars := mux.Vars(r)
	param := vars["questionnaireId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid questionnaire ID")
		return
	}

	questionnaire, err := app.models.Questionnaires.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	if userID != questionnaire.UserId {
		app.respondWithError(w, http.StatusInternalServerError, "Cannot delete other's questionnaire")
		return
	}

	err = app.models.Questionnaires.Delete(id)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
	}

	app.respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}
