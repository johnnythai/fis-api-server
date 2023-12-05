package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
	"github.com/julienschmidt/httprouter"
)

const version = "1.0.0"

type config struct {
	port int
	env string
}

type application struct {
	config config
	logger *slog.Logger
}

func main () {
	var cfg config
	
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.Parse()

	// declare logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// declare server
	app := &application{
		config: cfg,
		logger: logger,
	}

	// declare router	
	router := httprouter.New()

	// declare routes
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	
	// declare server
	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", cfg.port),
		Handler: router,
		IdleTimeout: time.Minute,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 30 * time.Second,
		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("Starting %s server on %s", cfg.env, srv.Addr)

	err := srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)	
}
