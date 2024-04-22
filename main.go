package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	// health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "success",
		})
	})

	handler := NewHandler()
	router.GET("/hotels", handler.GetHotels)

	// create and run the server
	s := &http.Server{
		Addr:              "localhost:8080",
		Handler:           router,
		ReadHeaderTimeout: 1 * time.Second,
		WriteTimeout:      3 * time.Second,
	}

	s.ListenAndServe()
}
