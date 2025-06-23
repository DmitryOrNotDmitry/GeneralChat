package repo

import (
	"generalChat/internal/model"

	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type ChatCache struct {
	redis *redis.Client
}

func CreateChatCache() *ChatCache {
	return &ChatCache{redis.NewClient(&redis.Options{
		Addr: "redis-cache:6379",
	})}
}

const chatKey = "chat:messages"
const maxMessages = 20

func (c *ChatCache) AddMessage(msg model.Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := c.redis.LPush(ctx, chatKey, data).Err(); err != nil {
		return err
	}

	return c.redis.LTrim(ctx, chatKey, 0, maxMessages-1).Err()
}

func (c *ChatCache) GetRecentMessages() ([]model.Message, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	values, err := c.redis.LRange(ctx, chatKey, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	messages := make([]model.Message, 0, len(values))
	for _, val := range values {
		var msg model.Message
		if err := json.Unmarshal([]byte(val), &msg); err == nil {
			messages = append(messages, msg)
		}
	}
	return messages, nil
}

func (c *ChatCache) Close() error {
	return c.redis.Close()
}
