package database

import (
	"context"
	"time"
)

type UserAuth struct {
	UserID         string    `db:"user_id"`
	Created        time.Time `db:"created"`
	HashedPassword string    `db:"hashed_password"`
}

func (db *DB) InsertToUserAuth(ctx context.Context, userID, hashedPassword string) error {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	query := `
		INSERT INTO user_auth (user_id, created, hashed_password)
		VALUES ($1, $2, $3)`

	_, err := db.ExecContext(ctx, query, userID, time.Now(), hashedPassword)
	if err != nil {
		return err
	}

	return err
}

func (db *DB) GetUserAuthByUserID(userID string) (*UserAuth, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	var user UserAuth

	query := `SELECT * FROM user_auth WHERE user_id = $1`

	err := db.GetContext(ctx, &user, query, userID)

	return &user, err
}

func (db *DB) UpdateUserAuthHashedPassword(userID, hashedPassword string) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `UPDATE user_auth SET hashed_password = $1 WHERE user_id = $2`

	_, err := db.ExecContext(ctx, query, hashedPassword, userID)
	return err
}
