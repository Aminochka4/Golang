package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"

	"github.com/Aminochka4/Golang/final-project/pkg/my-project/model"
	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

type config struct {
	port string
	env  string
	db   struct {
		dsn string
	}
}

type application struct {
	config config
	models model.Models
}

func main() {
	var cfg config
	flag.StringVar(&cfg.port, "port", ":8081", "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://postgres:Idinahui12345@localhost/postgres?sslmode=disable", "PostgreSQL DSN")
	flag.Parse()

	// Connect to DB
	db, err := openDB(cfg)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()

	app := &application{
		config: cfg,
		models: model.NewModels(db),
	}

	app.run()
}

func (app *application) run() {
	log.Println("Starting API server")

	r := mux.NewRouter()

	v1 := r.PathPrefix("/api/v1").Subrouter()

	//user

	v1.HandleFunc("/users/filter", app.getUsersByNameHandler).Methods("GET")

	v1.HandleFunc("/users/sorted", app.getUsersBySurnameHandler).Methods("GET")

	v1.HandleFunc("/users/limit", app.getUsersWithPaginationHandler).Methods("GET")

	v1.HandleFunc("users/register", app.registerUserHandler).Methods("POST")

	v1.HandleFunc("/users", app.getAllUsersHandler).Methods("GET")

	v1.HandleFunc("/users/{userId:[0-9]+}", app.getUserByIdHandler).Methods("GET")

	v1.HandleFunc("/users/{userId:[0-9]+}", app.updateUserHandler).Methods("PUT")

	v1.HandleFunc("/users/{userId:[0-9]+}", app.deleteUserHandler).Methods("DELETE")

	//questionnaire

	v1.HandleFunc("/questionnaire", app.createQuestionnaireHandler).Methods("POST")

	v1.HandleFunc("/questionnaire/{questionnaireId:[0-9]+}", app.getQuestionnaireHandler).Methods("GET")

	v1.HandleFunc("/questionnaire/{questionnaireId:[0-9]+}", app.updateQuestionnaireHandler).Methods("PUT")

	v1.HandleFunc("/questionnaire/{questionnaireId:[0-9]+}", app.deleteQuestionnaireHandler).Methods("DELETE")

	http.ListenAndServe(":8081", r)

	log.Printf("Starting server on %s\n", app.config.port)
	err := http.ListenAndServe(app.config.port, r)
	log.Fatal(err)
}

func openDB(cfg config) (*sql.DB, error) {
	// Use sql.Open() to create an empty connection pool, using the DSN from the config // struct.
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
