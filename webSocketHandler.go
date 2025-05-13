package main

import (
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
		var data map[string]string
		err := ws.ReadJSON(&data)
		if err != nil {
			delete(clients, ws)
			break
		}

		broadcastMessage(data)

		chatdb.SaveMessage(gin.H{
			"username":  data["username"],
			"message":   data["message"],
			"timestamp": time.Now(),
		})
	}
}

func broadcastMessage(msg map[string]string) {
	for client := range clients {
		err := client.WriteJSON(msg)
		if err != nil {
			client.Close()
			delete(clients, client)
		}
	}
}
