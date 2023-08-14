package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"example.com/assets"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/lib/pq"
)

type dbConnType int

const defaultTimeout = 3 * time.Second
const dbConnKey dbConnType = iota

type DB struct {
	*sqlx.DB
}

type Transactional interface {
	Commit() error
	Rollback() error
}

func New(dsn string, automigrate bool) (*DB, error) {
	db, err := sqlx.Connect("postgres", "postgres://"+dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(2 * time.Hour)

	if automigrate {
		iofsDriver, err := iofs.New(assets.EmbeddedFiles, "migrations")
		if err != nil {
			return nil, err
		}

		migrator, err := migrate.NewWithSourceInstance("iofs", iofsDriver, "postgres://"+dsn)
		if err != nil {
			return nil, err
		}

		err = migrator.Up()
		switch {
		case errors.Is(err, migrate.ErrNoChange):
			break
		case err != nil:
			return nil, err
		}
	}

	return &DB{db}, nil
}
func (db *DB) SetTxToContext(ctx context.Context, tx Transactional) context.Context {
	return context.WithValue(ctx, dbConnKey, tx)
}

func (db *DB) Transactional(ctx context.Context, f func(ctx context.Context) error) error {
	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("%w: fail to open transaction", err)
	}

	ctx = db.SetTxToContext(ctx, tx)
	if err := f(ctx); err != nil {
		if txErr := tx.Rollback(); txErr != nil {
			return fmt.Errorf("%w: fail to rollback transaction", txErr)
		}
		return err
	}

	if txErr := tx.Commit(); txErr != nil {
		return fmt.Errorf("%w: fail to commit transaction", txErr)
	}

	return nil
}
