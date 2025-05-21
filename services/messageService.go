package services

import (
	"fmt"
	"generalChat/entity"
)

type MessageRepository interface {
	SaveMessage(message entity.Message)
	GetLastNMessages(n int64) []entity.Message
}

type MessageCache interface {
	AddMessage(msg entity.Message) error
	GetRecentMessages() ([]entity.Message, error)
}

type MessageService struct {
	ChatRepo  MessageRepository
	ChatCache MessageCache
}

func (ms *MessageService) GetLast20Messages() []entity.Message {
	cachesMsgs, _ := ms.ChatCache.GetRecentMessages()
	if len(cachesMsgs) < 20 {
		fmt.Println("Данные берутся из БД")
		msgs := ms.ChatRepo.GetLastNMessages(20)

		for i := len(msgs) - 1; i >= 0; i-- {
			ms.ChatCache.AddMessage(msgs[i])
		}

		return msgs
	}

	fmt.Println("Данные берутся из кэша")
	return cachesMsgs
}
