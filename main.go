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
	// parse the configuration file
	cfg, err := LoadConfig("config.yaml")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration file")
	}

	// set up the service layer with the suppliers registry
	hotelService := NewUsecaseImpl(setupSupplierRegistry(cfg))
	// set up the handler layer
	handler := NewHandler(hotelService)
	// set up the router
	router := gin.Default()
	router.GET("/hotels", handler.GetHotels)

	// set up health check
	// health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "success",
		})
	})

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

// setupSupplierRegistry sets up the supplier registry based on the configuration.
func setupSupplierRegistry(cfg *Config) map[string]HotelSupplier {
	suppliers := make(map[string]HotelSupplier)

	// initialise specific supplier instances based on the configuration
	for name, sCfg := range cfg.Suppliers {
		switch name {
		case "Acme":
			suppliers[name] = supplier.NewAcme(sCfg.URL)
		case "Patagonia":
			suppliers[name] = supplier.NewPatagonia(sCfg.URL)
		case "Paperflies":
			suppliers[name] = supplier.NewPaperflies(sCfg.URL)
		default:
			log.Warn().Str("supplier", name).Msg("Unknown supplier, skipping")
		}
	}

	return suppliers
}
