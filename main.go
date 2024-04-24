package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ascenda-hotel/supplier"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func main() {
	router := gin.Default()
	// health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "success",
		})
	})

	// set up the service layer with the suppliers registry
	hotelService := NewUsecaseImpl(map[string]HotelSupplier{
		"Acme":       supplier.NewAcme("https://5f2be0b4ffc88500167b85a0.mockapi.io/suppliers/acme"),
		"Patagonia":  supplier.NewPatagonia("https://5f2be0b4ffc88500167b85a0.mockapi.io/suppliers/patagonia"),
		"Paperflies": supplier.NewPaperflies("https://5f2be0b4ffc88500167b85a0.mockapi.io/suppliers/paperflies"),
	})
	// set up the handler layer
	handler := NewHandler(hotelService)
	// set up the router
	router.GET("/hotels", handler.GetHotels)

	// create the server
	s := &http.Server{
		Addr:              "localhost:8080",
		Handler:           router,
		ReadHeaderTimeout: 2 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	// run the server
	go func() {
		// service connections
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	// set up graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("Shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Failed to shutdown server")
	}

	select {
	case <-ctx.Done():
		log.Info().Msg("Server shutdown timed out")
	}

	log.Info().Msg("Server shutdown complete")
}
