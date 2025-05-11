package main

import (
	"html/template"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()
	r.Use(cors.Default())

	r.SetHTMLTemplate(template.Must(template.ParseFiles("resources/index.html")))

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	r.GET("/api/data", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, CORS!",
			"status":  "ok",
		})
	})

	r.GET("/ws", handleConnections)

	log.Println("Сервер запущен на :8080")
	r.Run(":8080")
}
