package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// routes is our main application's router.
func (app *application) routes() http.Handler {
	r := mux.NewRouter()
	log.Println("Starting API server")
	// Convert the app.notFoundResponse helper to a http.Handler using the http.HandlerFunc()
	// adapter, and then set it as the custom error handler for 404 Not Found responses.
	r.NotFoundHandler = http.HandlerFunc(app.notFoundResponse)

	// Convert app.methodNotAllowedResponse helper to a http.Handler and set it as the custom
	// error handler for 405 Method Not Allowed responses
	r.MethodNotAllowedHandler = http.HandlerFunc(app.methodNotAllowedResponse)

	r.HandleFunc("/api/v1/healthcheck", app.healthcheckHandler).Methods("GET")

	//user

	user1 := r.PathPrefix("/api/v1").Subrouter()

	user1.HandleFunc("/users/register", app.registerUserHandler).Methods("POST")

	user1.HandleFunc("/users/activated", app.activateUserHandler).Methods("PUT")

	user1.HandleFunc("/users/login", app.createAuthenticationTokenHandler).Methods("POST")

	user1.HandleFunc("/users", app.getAllUsersHandler).Methods("GET")

	user1.HandleFunc("/users/{userId:[0-9]+}", app.getUserByIdHandler).Methods("GET")

	//questionnaire

	questionnaire1 := r.PathPrefix("/api/v1").Subrouter()

	questionnaire1.HandleFunc("/questionnaire", app.createQuestionnaireHandler).Methods("POST")

	questionnaire1.HandleFunc("/questionnaire", app.getAllQuestionnairesHandler).Methods("GET")

	questionnaire1.HandleFunc("/questionnaire/{questionnaireId:[0-9]+}", app.getQuestionnaireHandler).Methods("GET")

	questionnaire1.HandleFunc("/questionnaire/{questionnaireId:[0-9]+}", app.updateQuestionnaireHandler).Methods("PUT")

	questionnaire1.HandleFunc("/questionnaire/{questionnaireId:[0-9]+}", app.deleteQuestionnaireHandler).Methods("DELETE")

	//answer

	answer1 := r.PathPrefix("/api/v1").Subrouter()

	answer1.HandleFunc("/answer", app.createAnswerHandler).Methods("POST")

	answer1.HandleFunc("/answer", app.getAllAnswersHandler).Methods("GET")

	answer1.HandleFunc("/answer/{answerId:[0-9]+}", app.getAnswerHandler).Methods("GET")

	answer1.HandleFunc("/answer/{answerId:[0-9]+}", app.updateAnswerHandler).Methods("PUT")

	answer1.HandleFunc("/answer/{answerId:[0-9]+}", app.deleteAnswerHandler).Methods("DELETE")

	answer1.HandleFunc("/questionnaire/{questionnaireId:[0-9]+}/answer", app.getAnswerByQuestionnaireHandler).Methods("GET")

	return app.authenticate(r)
}
