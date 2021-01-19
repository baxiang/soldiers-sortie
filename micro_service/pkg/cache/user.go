package cache

import (
	"time"

	"github.com/go-redis/redis"
)

func GetUserInfo(client *redis.Client, userId string) (string, error) {
	return client.Get(userId).Result()
}

func SetUserInfo(client *redis.Client, userId string, value string, expiredAt time.Duration) error {
	return client.Set(userId, value, expiredAt).Err()
}

func DelUserInfo(client *redis.Client, userId string) error {
	return client.Del(userId).Err()
}
