package main

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"ascenda-hotel/entity"
)

const (
	// ErrInvalidDestinationID is returned when the destination ID is invalid.
	ErrInvalidDestinationID = "Invalid destination ID. Destination ID must be a non-negative integer."
	// ErrInternalServerError is returned when an internal server error occurs.
	ErrInternalServerError = "Internal server error. Please try again later or contact support."
	// ErrNoHotelsFound is returned when no hotels are found.
	ErrNoHotelsFound = "No hotels found."
)

type Usecase interface {
	// GetHotels returns a slice of Hotels for the given hotelIDs and destinationID.
	// If you do not want to filter by destinationID, set it to -1.
	// If both are provided, only hotels for the destinationID are returned.
	// If neither are provided, all hotels are returned.
	// If there is no matching hotel, it returns an empty slice.
	GetHotels(ctx context.Context, hotelIDs []string, destinationID int) ([]entity.Hotel, error)
}

type Handler struct {
	hotelService Usecase
}

func NewHandler(hotels Usecase) *Handler {
	return &Handler{
		hotelService: hotels,
	}
}

func (h *Handler) GetHotels(c *gin.Context) {
	// there are two optional query parameters: ids and destination_id
	// hotels is a comma-separated list of hotel IDs to retrieve
	// destination is the ID of the destination to retrieve hotels for
	// if both are provided, only hotels for the destination_id are returned
	// if neither are provided, all hotels are returned

	// parse the query params
	hotels := c.Query("hotels")
	destination := c.Query("destination")

	// if ids is provided, parse it into a slice of hotel IDs
	var hotelIDs []string
	if hotels != "" {
		hotelIDs = strings.Split(hotels, ",")
	}

	// parse the destination ID
	var destinationID int
	var err error
	if destination != "" {
		destinationID, err = strconv.Atoi(destination)
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": ErrInvalidDestinationID})
			return
		}
	} else {
		destinationID = -1
	}

	results, err := h.hotelService.GetHotels(c, hotelIDs, destinationID)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": ErrInternalServerError})
		return
	}

	if len(results) == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": ErrNoHotelsFound})
		return
	}

	c.JSON(http.StatusOK, results)
}
