package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"go-api/internal/data"
	"go-api/internal/driver"
)

type config struct {
	port int // what port do we want the web server to listen on
}

type application struct {
	config      config
	infoLog     *log.Logger
	errorLog    *log.Logger
	models      data.Models
	environment string
}

func main() {
	var cfg config
	cfg.port = 8082

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	dsn := os.Getenv("DSN")
	environment := os.Getenv("ENV")

	db, err := driver.ConnectPostgres(dsn)
	if err != nil {
		log.Fatal("Cannot connect to database")
	}
	defer db.SQL.Close()

	app := &application{
		config:      cfg,
		infoLog:     infoLog,
		errorLog:    errorLog,
		models:      data.New(db.SQL),
		environment: environment,
	}

	err = app.serve()
	if err != nil {
		log.Fatal(err)
	}
}

// serve starts the web server
func (app *application) serve() error {
	app.infoLog.Println("API listening on port", app.config.port)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.config.port),
		Handler: app.routes(),
	}

	return srv.ListenAndServe()
}
