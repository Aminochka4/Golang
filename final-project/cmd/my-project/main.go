package main

import (
	"database/sql"
	"flag"
	"github.com/Aminochka4/Golang/final-project/pkg/jsonlog"
	"github.com/Aminochka4/Golang/final-project/pkg/my-project/model/filler"
	"github.com/Aminochka4/Golang/final-project/pkg/vcs"
	"os"
	"sync"

	"github.com/Aminochka4/Golang/final-project/pkg/my-project/model"
	_ "github.com/lib/pq"
)

var (
	version = vcs.Version()
)

type config struct {
	port int
	env  string
	fill bool
	db   struct {
		dsn string
	}
}

type application struct {
	config config
	models model.Models
	logger *jsonlog.Logger
	wg     sync.WaitGroup
}

func main() {
	var cfg config
	flag.BoolVar(&cfg.fill, "fill", false, "Fill db with dummy data")
	flag.IntVar(&cfg.port, "port", 8081, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://postgres:Idinahui12345@localhost/postgres?sslmode=disable", "PostgreSQL DSN")
	flag.Parse()

	//Init logger
	logger := jsonlog.NewLogger(os.Stdout, jsonlog.LevelInfo)

	// Connect to DB
	db, err := openDB(cfg)
	if err != nil {
		logger.PrintError(err, nil)
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.PrintFatal(err, nil)
		}
	}()

	app := &application{
		config: cfg,
		models: model.NewModels(db),
		logger: logger,
	}

	if cfg.fill {
		err = filler.PopulateDatabase(app.models)
		if err != nil {
			logger.PrintFatal(err, nil)
			return
		}
	}

	if err := app.serve(); err != nil {
		logger.PrintFatal(err, nil)
	}
}

func openDB(cfg config) (*sql.DB, error) {
	// Use sql.Open() to create an empty connection pool, using the DSN from the config // struct.
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
