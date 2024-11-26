package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"shorter/internal/app/handlers"
)

func main() {
	//mux := http.NewServeMux()
	//mux.HandleFunc("/", handlers.URLHandler)
	//log.Fatal(http.ListenAndServe(":8080", mux))
	r := handlers.CreateRouter()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.POST("/", handlers.PostHandler)
	r.GET("/:id", handlers.GetHandler)
	log.Fatal(r.Run(":8080"))
}
