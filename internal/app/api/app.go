package api

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/0vkanix/greenlight/internal/movie"
)

const version = "1.0.0"

type Application struct {
	Config Config
	Logger *slog.Logger
	Movies movie.RepositoryInterface
}

func Run(cfg Config, logger *slog.Logger) error {
	conn, err := openDB(cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	logger.Info("database connection pool established")

	app := &Application{
		Config: cfg,
		Logger: logger,
		Movies: movie.NewRepository(conn),
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	if cfg.Env == "test" {
		return nil
	}

	logger.Info("starting server", "addr", srv.Addr, "env", cfg.Env)

	return srv.ListenAndServe()
}
