package main

import (
	"log/slog"
	"os"

	"github.com/0vkanix/greenlight/internal/app/api"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	cfg, err := api.ParseFlags(os.Args)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	if err := api.Run(cfg, logger); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
