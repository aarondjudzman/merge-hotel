package main

import (
	"context"
	"strings"

	"ascenda-hotel/entity"

	"github.com/rs/zerolog/log"
	"github.com/sourcegraph/conc/pool"
)

// HotelSupplier is an interface that defines the methods for fetching hotel data from a supplier.
type HotelSupplier interface {
	// FetchHotels returns a slice of Hotels for the given hotelIDs and destinationID.
	// If you do not want to filter by destinationID, set it to -1.
	// If both are provided, only hotels for the destinationID are returned.
	// If neither are provided, all hotels are returned.
	// If there is no matching hotel, it returns an empty slice.
	FetchHotels(ctx context.Context, hotelIDs []string, destinationID int) ([]entity.Hotel, error)
	GetName() string
}

// UsecaseImpl is a concrete implementation of the Usecase interface.
type UsecaseImpl struct {
	supplierRegistry map[string]HotelSupplier
}

// NewUsecaseImpl creates a new instance of the UsecaseImpl struct with the given supplierRegistry.
func NewUsecaseImpl(supplierRegistry map[string]HotelSupplier) *UsecaseImpl {
	return &UsecaseImpl{
		supplierRegistry: supplierRegistry,
	}
}

func (u *UsecaseImpl) GetHotels(ctx context.Context, hotelIDs []string, destinationID int) ([]entity.Hotel, error) {
	// concurrently fetch data from all suppliers
	p := pool.NewWithResults[[]entity.Hotel]()
	for _, supplier := range u.supplierRegistry {
		supplier := supplier // capture the loop variable
		p.Go(func() []entity.Hotel {
			logger := log.With().Str("supplier", supplier.GetName()).Logger()
			logger.Debug().Msgf("Fetching hotels from supplier %s", supplier.GetName())
			supplierHotels, err := supplier.FetchHotels(ctx, hotelIDs, destinationID)
			if err != nil {
				// if there is any error when fetching hotels from a supplier, we log it and return an empty slice
				// we do not return an error here because we want to continue fetching hotels from other suppliers
				logger.Error().Err(err).Msgf("Failed to fetch hotels from supplier %s", supplier.GetName())
				return []entity.Hotel{}
			}

			// clean the hotel data before returning it
			// doing this in service layer so that all the suppliers can use the same cleaner
			return cleanHotelData(supplierHotels)
		})
	}
	results := p.Wait()

	// flatten the results into a single slice
	var allHotels []entity.Hotel
	for _, supplierHotels := range results {
		allHotels = append(allHotels, supplierHotels...)
	}

	return allHotels, nil
}

// cleanHotelData performs some basic cleaning on the hotel data before returning it to the caller.
func cleanHotelData(hotels []entity.Hotel) []entity.Hotel {
	for i, hotel := range hotels {
		// sanitize whitespace from hotel data
		hotels[i] = entity.TrimSpaceFromHotel(hotel)

		// for hotel amenities, each supplier may have a different naming style.
		// e.g. WiFi, Wifi, wifi, etc..
		// we need to normalize the names to a common format of all lowercase
		for i, amenity := range hotel.Amenities.General {
			hotel.Amenities.General[i] = strings.ToLower(amenity)
		}
		for i, amenity := range hotel.Amenities.Room {
			hotel.Amenities.Room[i] = strings.ToLower(amenity)
		}
	}

	return hotels
}
