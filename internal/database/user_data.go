package database

import (
	"context"
	"time"
)

type UserData struct {
	UserID     string    `db:"user_id"`
	FirstName  string    `db:"first_name"`
	SecondName string    `db:"second_name"`
	Birthdate  time.Time `db:"birthdate"`
	Sex        bool      `db:"sex"`
	Biography  string    `db:"biography"`
	City       string    `db:"city"`
}

func (db *DB) InsertToUserData(ctx context.Context, userData UserData) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	var userID string

	query, args, err := db.BindNamed(
		`INSERT INTO user_data (first_name, second_name, birthdate, sex, biography, city)
				VALUES (:first_name, :second_name, :birthdate, :sex, :biography, :city)
				RETURNING user_id`, &userData)

	if err != nil {
		return "", err
	}

	if err := db.GetContext(ctx, &userID, query, args...); err != nil {
		return "", err
	}

	return userID, nil
}

func (db *DB) GetUserDataByID(userID string) (*UserData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	var data UserData

	query := `SELECT *
		FROM user_data WHERE user_id = $1 `
	if err := db.GetContext(ctx, &data, query, userID); err != nil {
		return nil, err
	}

	return &data, nil
}
