package main

import (
	"encoding/json"
	"github.com/Aminochka4/Golang/final-project/pkg/my-project/model"
	"github.com/Aminochka4/Golang/final-project/pkg/my-project/validator"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

func (app *application) respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)

	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Surname  string `json:"surname"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
	}

	user := &model.User{
		Name:     input.Name,
		Surname:  input.Surname,
		Username: input.Username,
		Email:    input.Email,
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	v := validator.New()

	if model.ValidateUser(v, user); !v.Valid() {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	err = app.models.Users.Insert(user)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	err = app.models.Permissions.AddForUser(user.Id, "questionnaire:done")
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	token, err := app.models.Tokens.New(user.Id, 3*24*time.Hour, model.ScopeActivation)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	var res struct {
		Token *string     `json:"token"`
		User  *model.User `json:"user"`
	}

	res.Token = &token.Plaintext
	res.User = user

	app.writeJSON(w, http.StatusCreated, envelope{"user": res}, nil)
}

func (app *application) getAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := app.models.Users.GetAll()
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "Failed to fetch users")
		return
	}
	app.respondWithJson(w, http.StatusOK, users)
}

func (app *application) getUserByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["userId"]

	id, err := strconv.Atoi(param)

	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid user Id")
		return
	}

	user, err := app.models.Users.GetById(id)

	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}
	app.respondWithJson(w, http.StatusOK, user)
}

func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the plaintext activation token from the request body
	var input struct {
		TokenPlaintext string `json:"token"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	// Validate the plaintext token provided by the client.
	v := validator.New()

	if model.ValidateTokenPlaintext(v, input.TokenPlaintext); !v.Valid() {
		app.respondWithError(w, http.StatusBadRequest, "400 Invalid token plaintext")
		return
	}

	// Retrieve the details of the user associated with the token using the GetForToken() method.
	// If no matching record is found, then we let the client know that the token they provided
	// is not valid.
	user, err := app.models.Users.GetForToken(model.ScopeActivation, input.TokenPlaintext)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "400 Something wrong")
		return
	}

	user.Activated = true

	err = app.models.Users.Update(user)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "400 Something wrong")
		return
	}

	err = app.models.Tokens.DeleteAllForUser(model.ScopeActivation, user.Id)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "400 Something wrong")
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Status Internal Server Error")
	}
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		return err
	}

	return nil
}

// tsis3
func (app *application) getUsersByNameHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	if name == "" {
		app.respondWithError(w, http.StatusBadRequest, "Missing name parameter")
		return
	}

	users, err := app.models.Users.GetByName(name)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "Failed to fetch users")
		return
	}

	app.respondWithJSON(w, http.StatusOK, users)
}

func (app *application) getUsersBySurnameHandler(w http.ResponseWriter, r *http.Request) {
	users, err := app.models.Users.GetBySurname()
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "Failed to fetch users")
		return
	}

	app.respondWithJSON(w, http.StatusOK, users)
}

func (app *application) getUsersWithPaginationHandler(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid limit parameter")
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid offset parameter")
		return
	}

	users, err := app.models.Users.GetUsersWithPagination(limit, offset)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "Failed to fetch users")
		return
	}

	app.respondWithJSON(w, http.StatusOK, users)
}
