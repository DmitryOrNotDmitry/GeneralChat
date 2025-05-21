package infra

import (
	"context"
	"encoding/json"
	"generalChat/entity"
	"time"

	"github.com/redis/go-redis/v9"
)

type ChatCache struct {
	cache *redis.Client
}

func CreateChatCache() *ChatCache {
	return &ChatCache{redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})}
}

const chatKey = "chat:messages"
const maxMessages = 20

func (c *ChatCache) AddMessage(msg entity.Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := c.cache.LPush(ctx, chatKey, data).Err(); err != nil {
		return err
	}

	return c.cache.LTrim(ctx, chatKey, 0, maxMessages-1).Err()
}

func (c *ChatCache) GetRecentMessages() ([]entity.Message, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	values, err := c.cache.LRange(ctx, chatKey, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	messages := make([]entity.Message, 0, len(values))
	for _, val := range values {
		var msg entity.Message
		if err := json.Unmarshal([]byte(val), &msg); err == nil {
			messages = append(messages, msg)
		}
	}
	return messages, nil
}
