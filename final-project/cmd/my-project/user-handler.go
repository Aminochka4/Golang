package main

import (
	"encoding/json"
	"github.com/Aminochka4/Golang/final-project/pkg/my-project/model"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
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

func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
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
	hashed, _ := bcrypt.GenerateFromPassword([]byte(input.Password), 8)

	user := &model.User{
		Name:     input.Name,
		Surname:  input.Surname,
		Username: input.Username,
		Email:    input.Email,
		Password: string(hashed),
	}

	err = app.models.Users.Insert(user)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJson(w, http.StatusCreated, user)
}

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

func (app *application) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["userId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid menu ID")
		return
	}

	user, err := app.models.Users.GetById(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	var input struct {
		Name     *string `json:"name"`
		Surname  *string `json:"surname"`
		Username *string `json:"username"`
		Email    *string `json:"email"`
		Password *string `json:"password"`
	}

	err = app.readJSON(w, r, &input)

	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if input.Name != nil {
		user.Name = *input.Name
	}
	if input.Surname != nil {
		user.Surname = *input.Surname
	}
	if input.Username != nil {
		user.Username = *input.Username
	}
	if input.Email != nil {
		user.Email = *input.Email
	}
	if input.Password != nil {
		user.Password = *input.Password
	}

	err = app.models.Users.Update(user)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
	}

	app.respondWithJson(w, http.StatusOK, user)
}

func (app *application) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	params := vars["userId"]

	id, err := strconv.Atoi(params)

	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid User Id")
		return
	}

	err = app.models.Users.Delete(id)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
	}

	app.respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
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
