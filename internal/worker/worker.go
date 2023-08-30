package worker

import (
	"context"
	"encoding/json"

	"otus-homework/internal/database"
	"otus-homework/internal/redis"
)

type Worker struct {
	Cache *redis.Cache
	DB    *database.DB
}

const AddPost = "addPost"

func New(cache *redis.Cache, db *database.DB) Worker {
	return Worker{
		Cache: cache,
		DB:    db,
	}
}

func (w *Worker) GetTaskFromQueue() (int64, error) {
	return w.DB.GetTaskForProcess(context.Background())
}

func (w *Worker) ProcessTask(id int64) error {
	task, err := w.DB.GetTaskByID(context.Background(), id)
	if err != nil {
		return err
	}
	var payload = database.Payload{}
	err = json.Unmarshal([]byte(task.Payload), &payload)
	if err != nil {
		return err
	}
	friendIDs, err := w.DB.GetUserFriendsIDsByUserID(payload.UserID)
	if err != nil {
		return err
	}

	for _, friendID := range friendIDs {
		w.Cache.Put(friendID, payload.Post, 1000)
	}

	return w.DB.SetTaskInQueue(context.Background(), id)
}

func (w *Worker) AddTasksToQueue(payload string) error {
	return w.DB.CreateTask(context.Background(), AddPost, payload)
}
