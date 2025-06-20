package service

import (
	"generalChat/internal/model"
	"log"
)

type MessageRepository interface {
	SaveMessage(message model.Message) error
	GetLastNMessages(n int64) ([]model.Message, error)
}

type MessageCache interface {
	AddMessage(msg model.Message) error
	GetRecentMessages() ([]model.Message, error)
}

type MessageService struct {
	ChatRepo  MessageRepository
	ChatCache MessageCache
}

func (ms *MessageService) GetLast20Messages() []model.Message {
	cachesMsgs, _ := ms.ChatCache.GetRecentMessages()
	if len(cachesMsgs) < 20 {
		log.Println("Данные берутся из БД")
		msgs, _ := ms.ChatRepo.GetLastNMessages(20)

		for i := len(msgs) - 1; i >= 0; i-- {
			ms.ChatCache.AddMessage(msgs[i])
		}

		return msgs
	}

	log.Println("Данные берутся из кэша")
	return cachesMsgs
}
