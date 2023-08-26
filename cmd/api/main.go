package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/debug"
	"sync"

	"example.com/internal/database"
	"example.com/internal/env"
	"example.com/internal/version"

	"golang.org/x/exp/slog"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	err := run(logger)
	if err != nil {
		trace := string(debug.Stack())
		logger.Error(err.Error(), slog.String("trace", trace))
		os.Exit(1)
	}
}

type config struct {
	baseURL   string
	httpPort  int
	basicAuth struct {
		username       string
		hashedPassword string
	}
	db struct {
		dsn         string
		automigrate bool
	}
	dbReplicaDSN string
	jwt          struct {
		secretKey string
	}
}

type application struct {
	config    config
	db        *database.DB
	dbReplica *database.DB
	logger    *slog.Logger
	wg        sync.WaitGroup
}

func run(logger *slog.Logger) error {
	var cfg config

	cfg.baseURL = env.GetString("BASE_URL", "http://localhost:4444")
	cfg.httpPort = env.GetInt("HTTP_PORT", 4444)
	cfg.basicAuth.username = env.GetString("BASIC_AUTH_USERNAME", "admin")
	cfg.basicAuth.hashedPassword = env.GetString("BASIC_AUTH_HASHED_PASSWORD", "$2a$10$jRb2qniNcoCyQM23T59RfeEQUbgdAXfR6S0scynmKfJa5Gj3arGJa")
	cfg.db.dsn = env.GetString("DB_DSN", "postgres:postgres@localhost:5432/postgres?sslmode=disable")
	cfg.dbReplicaDSN = env.GetString("DB_REPLICA_DSN", "postgres:postgres@localhost:5433/postgres?sslmode=disable")
	cfg.db.automigrate = env.GetBool("DB_AUTOMIGRATE", true)
	cfg.jwt.secretKey = env.GetString("JWT_SECRET_KEY", "iugey2xd4ctpeaefpnmy3nuvzj6ewsm3")

	showVersion := flag.Bool("version", false, "display version and exit")

	flag.Parse()

	if *showVersion {
		fmt.Printf("version: %s\n", version.Get())
		return nil
	}

	db, err := database.New(cfg.db.dsn, cfg.db.automigrate)
	if err != nil {
		return err
	}
	defer db.Close()
	dbReplica, err := database.New(cfg.dbReplicaDSN, false)
	if err != nil {
		return err
	}
	defer dbReplica.Close()

	app := &application{
		config:    cfg,
		db:        db,
		dbReplica: dbReplica,
		logger:    logger,
	}

	return app.serveHTTP()
}
