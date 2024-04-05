package main

import (
	"encoding/json"
	"github.com/Aminochka4/Golang/final-project/pkg/my-project/model"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (app *application) respondWithError(w http.ResponseWriter, code int, message string) {
	app.respondWithJSON(w, code, map[string]string{"error": message})
}

func (app *application) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)

	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (app *application) createQuestionnaireHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Topic     string `json:"topic"`
		Questions string `json:"questions"`
		UserId    int64  `json:"userId"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
	}

	questionnaire := &model.Questionnaire{
		Topic:     input.Topic,
		Questions: input.Questions,
		UserId:    input.UserId,
	}

	err = app.models.Questionnaires.Insert(questionnaire)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusCreated, questionnaire)
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

	app.respondWithJSON(w, http.StatusOK, questionnaire)
}

func (app *application) updateQuestionnaireHandler(w http.ResponseWriter, r *http.Request) {
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

	app.respondWithJSON(w, http.StatusOK, questionnaire)
}

func (app *application) deleteQuestionnaireHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["questionnaireId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid questionnaire ID")
		return
	}

	err = app.models.Questionnaires.Delete(id)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
	}

	app.respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
