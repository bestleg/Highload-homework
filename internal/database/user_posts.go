package database

import (
	"context"
	"time"
)

type UserPost struct {
	ID      string    `db:"post_id"`
	UserID  string    `db:"user_id"`
	Post    string    `db:"post"`
	Created time.Time `db:"created"`
}

func (db *DB) InsertToUserPosts(ctx context.Context, userPost UserPost) error {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	query := `
		INSERT INTO user_posts (user_id, post, created)
		VALUES ($1, $2, $3)`

	_, err := db.ExecContext(ctx, query, userPost.UserID, userPost.Post, time.Now())

	return err
}

func (db *DB) UpdateUserPosts(ctx context.Context, userPost UserPost) error {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	query := `
		UPDATE user_posts SET 
		                       post = $1,
		                       created = $2
		    FROM user_posts WHERE post_id = $3
		`

	_, err := db.ExecContext(ctx, query, userPost.Post, time.Now(), userPost.ID)

	return err
}

func (db *DB) DeleteUserPosts(ctx context.Context, userPost UserPost) error {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	query := `
		DELETE FROM user_posts WHERE post_id = $1 AND user_id = $2
		`

	_, err := db.ExecContext(ctx, query, userPost.ID, userPost.UserID)

	return err
}

func (db *DB) GetUserPostByPostID(userID string) (*UserPost, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	var post UserPost

	query := `SELECT * FROM user_posts WHERE user_id = $1`

	err := db.GetContext(ctx, &post, query, userID)

	return &post, err
}
