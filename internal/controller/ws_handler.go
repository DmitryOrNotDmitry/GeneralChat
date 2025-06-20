package controller

import (
	"generalChat/internal/model"
	"generalChat/internal/service"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WSHandler struct {
	clients       map[*websocket.Conn]bool
	clientActions chan ClientAction
	mu            sync.RWMutex

	upgrader websocket.Upgrader

	chatdb    service.MessageRepository
	chatCache service.MessageCache
}

const (
	Join = iota
	Exit
)

type ClientAction struct {
	conn   *websocket.Conn
	action int
}

func ConstructorWSHandler(chatdb service.MessageRepository, chatCache service.MessageCache) *WSHandler {
	ws := &WSHandler{
		clients:       make(map[*websocket.Conn]bool),
		clientActions: make(chan ClientAction, 5),
		mu:            sync.RWMutex{},
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
		chatdb:    chatdb,
		chatCache: chatCache,
	}

	go func() {
		for client := range ws.clientActions {
			ws.mu.Lock()
			if client.action == Join {
				ws.clients[client.conn] = true
			} else if client.action == Exit {
				delete(ws.clients, client.conn)
			}
			ws.mu.Unlock()
		}
	}()

	return ws
}

func (wsh *WSHandler) HandleConnections(c *gin.Context) {

	ws, err := wsh.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	wsh.clientActions <- ClientAction{ws, Join}
	defer ws.Close()

	for {
		var msg model.Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			wsh.clientActions <- ClientAction{ws, Exit}
			break
		}
		msg.Timestamp = time.Now()

		wsh.broadcastMessage(msg)
		wsh.chatdb.SaveMessage(msg)
		wsh.chatCache.AddMessage(msg)
	}
}

func (wsh *WSHandler) broadcastMessage(msg model.Message) {
	actions := make([]ClientAction, 0)
	wsh.mu.RLock()
	for client := range wsh.clients {
		err := client.WriteJSON(msg)
		if err != nil {
			client.Close()
			actions = append(actions, ClientAction{client, Exit})
		}
	}
	wsh.mu.RUnlock()

	for _, action := range actions {
		wsh.clientActions <- action
	}
}

func (wsh *WSHandler) Close() {
	actions := make([]ClientAction, 0)
	wsh.mu.RLock()
	for client := range wsh.clients {
		client.Close()
		actions = append(actions, ClientAction{client, Exit})
	}
	wsh.mu.RUnlock()

	for _, action := range actions {
		wsh.clientActions <- action
	}

	close(wsh.clientActions)
}
