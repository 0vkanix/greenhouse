package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
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
	f.SetOutput(io.Discard)

	f.IntVar(&cfg.port, "port", 0, "API server port")
	f.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	f.StringVar(&cfg.db.dsn, "db-dsn", "", "PostgreSQL DSN")

	if err := f.Parse(args[1:]); err != nil {
		return err
	}

	if cfg.port == 0 {
		return fmt.Errorf("the -port flag is required")
	}
	if cfg.db.dsn == "" {
		return fmt.Errorf("the -db-dsn flag is required")
	}

	conn, err := openDB(cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	logger.Info("database connection pool established")

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

	if cfg.env == "test" {
		return nil
	}

	logger.Info("starting server", "addr", srv.Addr, "env", cfg.env)

	return srv.ListenAndServe()
}

func openDB(cfg config) (*pgxpool.Pool, error) {
	conn, err := pgxpool.New(context.Background(), cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = conn.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
