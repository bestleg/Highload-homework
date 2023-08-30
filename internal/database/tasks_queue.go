package database

import (
	"context"
	"database/sql"
	"fmt"
)

type Task struct {
	ID      int64  `db:"id"`
	Name    string `db:"name"`
	Payload string `db:"payload"`
	InQueue bool   `db:"in_queue"`
}

type Payload struct {
	UserID string `json:"user_id"`
	Post   string `json:"post"`
}

func (db *DB) CreateTask(ctx context.Context, name, payload string) error {

	query := `INSERT INTO tasks_queue (name, payload, in_queue) VALUES ($1, $2,false) ON CONFLICT DO NOTHING`
	if _, err := db.ExecContext(ctx, query, name, payload); err != nil {
		return err
	}

	return nil
}

func (db *DB) CreateTasks(ctx context.Context, tasks []Task) (int64, error) {
	query := `INSERT INTO tasks_queue(name, payload,in_queue)
		VALUES (:name, :payload,false)
		ON CONFLICT DO NOTHING `

	res, err := db.NamedExecContext(ctx, query, tasks)
	if err != nil {
		return 0, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return 0, sql.ErrNoRows
	}

	return affected, nil
}

func (db *DB) GetTaskByID(ctx context.Context, id int64) (*Task, error) {
	var task Task
	query := `SELECT id, name, payload FROM tasks_queue WHERE id = $1`
	if err := db.GetContext(ctx, &task, query, id); err != nil {
		return nil, fmt.Errorf("%w: tasks_queue.GetTaskByID fail", err)
	}

	return &task, nil
}

func (db *DB) GetTaskForProcess(ctx context.Context) (int64, error) {
	var task int64
	sql := `
		SELECT id
		FROM tasks_queue
		WHERE NOT in_queue
		ORDER BY id
		FOR UPDATE SKIP LOCKED
		LIMIT 1;
	`

	if err := db.GetContext(ctx, &task, sql); err != nil {
		return 0, fmt.Errorf("%w: tasks_queue.GetTaskForProcess fail", err)
	}

	return task, nil
}

func (db *DB) SetTaskInQueue(ctx context.Context, id int64) error {

	query := `
		UPDATE tasks_queue SET in_queue = true WHERE id = $1
		`

	if _, err := db.ExecContext(ctx, query, id); err != nil {
		return err
	}

	return nil
}
