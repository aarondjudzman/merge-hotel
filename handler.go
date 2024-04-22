package main

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
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

	fmt.Println("Hotels: ", hotelIDs)
	fmt.Println("Destination: ", destination)

	c.JSON(200, gin.H{
		"message": "Hotels retrieved successfully",
	})
}
