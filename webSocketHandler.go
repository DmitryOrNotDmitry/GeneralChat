package main

import (
	"generalChat/entity"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var clients = make(map[*websocket.Conn]bool)

func handleConnections(c *gin.Context) {

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()
	clients[ws] = true

	for {
		var msg entity.Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			delete(clients, ws)
			break
		}
		msg.Timestamp = time.Now()

		broadcastMessage(msg)
		chatdb.SaveMessage(msg)
		chatCache.AddMessage(msg)
	}
}

func broadcastMessage(msg entity.Message) {
	for client := range clients {
		err := client.WriteJSON(msg)
		if err != nil {
			client.Close()
			delete(clients, client)
		}
	}
}
