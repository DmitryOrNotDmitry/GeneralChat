package main

import "fmt"

func GetLast20Messages() []Message {
	cachesMsgs, _ := chatCache.GetRecentMessages()
	if len(cachesMsgs) < 20 {
		fmt.Println("Данные берутся из БД")
		msgs := chatdb.GetLastNMessages(20)

		for i := len(msgs) - 1; i >= 0; i-- {
			chatCache.AddMessage(msgs[i])
		}

		return msgs
	}

	fmt.Println("Данные берутся из Redis")
	return cachesMsgs
}
