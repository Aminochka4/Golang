package main

import (
	"net/http"
	"strings"

	"github.com/Aminochka4/Golang/final-project/pkg/my-project/model"
	"github.com/Aminochka4/Golang/final-project/pkg/my-project/validator"
)

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Vary", "Authorization")

		authorizationHeader := r.Header.Get("Authorization")

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			app.respondWithError(w, http.StatusInternalServerError, "500 Invalid authorization")
			return
		}

		token := headerParts[1]

		v := validator.New()

		if model.ValidateTokenPlaintext(v, token); !v.Valid() {
			app.respondWithError(w, http.StatusBadRequest, "400 Invalid authentication")
			return
		}

		user, err := app.models.Users.GetForToken(model.ScopeAuthentication, token)
		if err != nil {
			app.respondWithJson(w, http.StatusNotFound, "404 Invalid token")
			return
		}

		r = app.contextSetUser(r, user)

		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuthenticatedUser(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		next.ServeHTTP(w, r)
	}
}

func (app *application) requireActivatedUser(next http.HandlerFunc) http.HandlerFunc {
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)

		if !user.Activated {
			app.respondWithJSON(w, http.StatusInternalServerError, "500 Invalid activation")
			return
		}

		next.ServeHTTP(w, r)
	})

	return app.requireAuthenticatedUser(fn)
}

func (app *application) requirePermissions(code string, next http.HandlerFunc) http.HandlerFunc {
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)

		permissions, err := app.models.Permissions.GetAllForUser(user.Id)
		if err != nil {
			app.respondWithError(w, http.StatusBadRequest, "400 Invalid Permission")
			return
		}

		if !permissions.Include(code) {
			app.respondWithError(w, http.StatusInternalServerError, "500 Invalid permitted")
			return
		}

		next.ServeHTTP(w, r)
	})

	return app.requireActivatedUser(fn)
}
