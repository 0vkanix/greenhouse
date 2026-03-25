package api

import (
	"flag"
	"fmt"
	"time"
)

type Config struct {
	Port int
	Env  string
	DB   struct {
		DSN          string
		MaxOpenConns int
		MaxIdleConns int
		MaxIdleTime  time.Duration
	}
}

func ParseFlags(args []string) (Config, error) {
	var cfg Config

	f := flag.NewFlagSet(args[0], flag.ContinueOnError)

	f.IntVar(&cfg.Port, "port", 0, "API server port")
	f.StringVar(&cfg.Env, "env", "development", "Environment (development|staging|production)")
	f.StringVar(&cfg.DB.DSN, "db-dsn", "", "PostgreSQL DSN")
	f.IntVar(&cfg.DB.MaxOpenConns, "db-max-open-conns", 0, "PostgreSQL max open connections")
	f.IntVar(&cfg.DB.MaxIdleConns, "db-max-idle-conns", 0, "PostgreSQL max idle connections")
	f.DurationVar(&cfg.DB.MaxIdleTime, "db-max-idle-time", 0, "PostgreSQL max connection idle time")

	if err := f.Parse(args[1:]); err != nil {
		return Config{}, err
	}

	if cfg.Port == 0 {
		return Config{}, fmt.Errorf("the -port flag is required")
	}
	if cfg.DB.DSN == "" {
		return Config{}, fmt.Errorf("the -db-dsn flag is required")
	}

	return cfg, nil
}
