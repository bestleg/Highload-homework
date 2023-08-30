package database

import (
	"context"
	"database/sql"
	"time"
)

type UserFriend struct {
	UserID       string `db:"user_id"`
	FriendUserID string `db:"user_friend_id"`
}

func (db *DB) InsertToUserFriend(ctx context.Context, userFriend UserFriend) error {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	query := `
		INSERT INTO user_friend (user_id, user_friend_id, created)
		VALUES ($1, $2, $3)`

	_, err := db.ExecContext(ctx, query, userFriend.UserID, userFriend.FriendUserID, time.Now())

	return err
}

func (db *DB) DeleteUserFriendByUserID(userID, friendUserID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `DELETE FROM user_friend WHERE user_id = $1 and user_friend_id = $2`

	result, err := db.ExecContext(ctx, query, userID, friendUserID)

	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (db *DB) GetUserFriendsIDsByUserID(userID string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	var userFriendIDs []string

	query := `SELECT user_friend_id FROM user_friend WHERE user_id = $1`

	err := db.SelectContext(ctx, &userFriendIDs, query, userID)

	return userFriendIDs, err
}
