package database

import (
	"context"
)

type UserFeed []UserFeedSingle

type UserFeedSingle struct {
	UserID string `db:"user_id"`
	Post   string `db:"post"`
}

func (db *DB) GetAllFeed() (UserFeed, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	var posts []UserFeedSingle

	query := `SELECT user_friend.user_id, up.post
				FROM user_friend
    			JOIN user_posts as up ON user_friend.user_friend_id = up.user_id
				ORDER BY up.created`

	err := db.SelectContext(ctx, &posts, query)

	return posts, err
}
