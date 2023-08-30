package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) *Cache {
	return &Cache{
		client: client,
	}
}

func (rc *Cache) Put(key string, value string, limit int64) error {
	// Добавляем строку в конец массива по ключу
	err := rc.client.RPush(context.Background(), key, value).Err()
	if err != nil {
		return err
	}

	// Получаем количество элементов массива
	length, err := rc.client.LLen(context.Background(), key).Result()
	if err != nil {
		return err
	}

	// Если количество элементов превышает лимит, удаляем первый элемент
	if length > limit {
		err = rc.client.LPop(context.Background(), key).Err()
		if err != nil {
			return err
		}
	}

	return nil
}

func (rc *Cache) Get(key string) ([]string, error) {
	// Получаем все элементы массива по ключу
	values, err := rc.client.LRange(context.Background(), key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	return values, nil
}
