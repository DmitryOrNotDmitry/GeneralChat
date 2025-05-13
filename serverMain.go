package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var chatdb ChatDB = *CreateChatDB()

func main() {
	r := gin.Default()
	r.Use(cors.Default())

	r.Static("/resources/static", "./resources/static")

	r.SetHTMLTemplate(template.Must(template.ParseFiles("resources/index.html")))

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.GET("/api/lastMessages", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"messages": chatdb.GetLast20Messages(),
		})
	})

	r.GET("/ws", handleConnections)
	defer chatdb.Close()

	log.Println("Сервер запущен на :8080")
	r.Run(":8080")
}
