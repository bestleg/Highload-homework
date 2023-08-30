package main

import (
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"

	"otus-homework/internal/database"
	"otus-homework/internal/env"
	redisCache "otus-homework/internal/redis"
	"otus-homework/internal/worker"
)

func main() {
	dsn := env.GetString("DB_DSN", "postgres:postgres@localhost:5432/postgres?sslmode=disable")

	db, err := database.New(dsn, false)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()
	cache := redisCache.NewRedisCache(redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	}))

	worker := worker.New(cache, db)

	// греем холодный кеш, так конечно нельзя, но в рамках теста можно :)
	feed, err := db.GetAllFeed()
	if err != nil {
		log.Panic(err)
	}

	for _, single := range feed {
		err = cache.Put(single.UserID, single.Post, 1000)
		if err != nil {
			fmt.Println("Ошибка redis:", err)
		}
		err = nil
	}

	for {
		// Получаем задачу из очереди
		task, err := worker.GetTaskFromQueue()
		if err != nil {
			fmt.Println("Ошибка получения задачи из очереди:", err)
			time.Sleep(time.Second)
			continue
		}

		// Обрабатываем задачу
		err = worker.ProcessTask(task)
		if err != nil {
			fmt.Println("Ошибка обработки задачи:", err)
		}
	}

}
