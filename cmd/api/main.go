package main

import (
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
}

type application struct {
	config config
	logger *slog.Logger
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	if err := run(os.Args, os.Stdout, logger); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func run(args []string, stdout io.Writer, logger *slog.Logger) error {
	var cfg config

	f := flag.NewFlagSet(args[0], flag.ContinueOnError)
	f.SetOutput(io.Discard) // keep tests quiet

	f.IntVar(&cfg.port, "port", 4000, "API server port")
	f.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	if err := f.Parse(args[1:]); err != nil {
		return err
	}

	app := &application{
		config: cfg,
		logger: logger,
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	// For testing purposes, we can check for a special environment "test-config"
	// to avoid starting the server.
	if cfg.env == "test-config" {
		return nil
	}

	logger.Info("starting server", "addr", srv.Addr, "env", cfg.env)

	return srv.ListenAndServe()
}
