package main

import (
	"github.com/Aminochka4/Golang/final-project/pkg/my-project/model"
	"github.com/gorilla/mux"
	"strconv"

	//"github.com/gorilla/mux"
	"log"
	"net/http"
	//"strconv"
)

func (app *application) getAllAnswersHandler(w http.ResponseWriter, r *http.Request) {
	answer, err := app.models.Answer.GetAll()
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "Failed to fetch questionnaires")
		return
	}
	app.respondWithJson(w, http.StatusOK, answer)
}

func (app *application) createAnswerHandler(w http.ResponseWriter, r *http.Request) {
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
		QuestionnaireId string `json:"questionnaireId"`
		Answer          string `json:"answer"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		log.Println(err)
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}

	answer := &model.Answer{
		QuestionnaireId: input.QuestionnaireId,
		Answer:          input.Answer,
		UserId:          userID,
	}

	err = app.models.Answer.Insert(answer)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJson(w, http.StatusCreated, answer)
}

func (app *application) getAnswerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["answerId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid answer ID")
		return
	}

	answer, err := app.models.Answer.Get(id)
	//log.Println("There is a error")
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	app.respondWithJson(w, http.StatusOK, answer)
}

func (app *application) updateAnswerHandler(w http.ResponseWriter, r *http.Request) {
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
	param := vars["answerId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid answer ID")
		return
	}

	answer, err := app.models.Answer.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	var input struct {
		QuestionnaireId *string `json:"questionnaireId"`
		Answer          *string `json:"answer"`
	}

	if userID != answer.UserId {
		app.respondWithError(w, http.StatusInternalServerError, "Cannot update other's answer")
		return
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if input.QuestionnaireId != nil {
		answer.QuestionnaireId = *input.QuestionnaireId
	}

	if input.Answer != nil {
		answer.Answer = *input.Answer
	}

	err = app.models.Answer.Update(answer)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
	}

	app.respondWithJson(w, http.StatusOK, answer)
}

func (app *application) deleteAnswerHandler(w http.ResponseWriter, r *http.Request) {
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
	param := vars["answerId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid answer ID")
		return
	}

	answer, err := app.models.Answer.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	if userID != answer.UserId {
		app.respondWithError(w, http.StatusInternalServerError, "Cannot delete other's answer")
		return
	}

	err = app.models.Answer.Delete(id)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
	}

	app.respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

func (app *application) getAnswerByQuestionnaireHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["questionnaireId"] // Обратите внимание на использование "questionnaireId" здесь

	questionnaireID, err := strconv.Atoi(param)
	if err != nil || questionnaireID < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid questionnaire ID")
		return
	}

	answers, err := app.models.Answer.GetByQuestionnaire(questionnaireID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.respondWithJson(w, http.StatusOK, answers)
}
