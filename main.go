package main

import (
	"generalChat/internal/controller"
	"generalChat/internal/repo"
	"generalChat/internal/service"

	"html/template"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var chatdb = repo.CreateChatDB()
var chatCache = repo.CreateChatCache()
var wshandler = controller.ConstructorWSHandler(chatdb, chatCache)

var chatService = &service.MessageService{
	ChatRepo:  chatdb,
	ChatCache: chatCache,
}

func main() {
	defer chatdb.Close()
	defer wshandler.Close()
	defer chatCache.Close()

	controller.InitPrometheus()

	r := gin.Default()
	r.Use(cors.Default())
	r.Use(controller.MetricsMiddleware())
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	r.Static("resources/static", "resources/static")

	r.SetHTMLTemplate(template.Must(template.ParseFiles("resources/index.html")))

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.GET("/api/lastMessages", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"messages": chatService.GetLast20Messages(),
		})
	})

	r.GET("/ws", wshandler.HandleConnections)

	log.Println("Сервер запущен на :8080")
	r.Run(":8080")
}
