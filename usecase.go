package main

import (
	"context"
	"errors"
	"strings"
	"time"

	"merge-hotel/entity"

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

type Cacher interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{}, ttl time.Duration)
}

// UsecaseImpl is a concrete implementation of the Usecase interface.
type UsecaseImpl struct {
	supplierRegistry map[string]HotelSupplier
	cache            Cacher
}

// NewUsecaseImpl creates a new instance of the UsecaseImpl struct with the given supplierRegistry.
func NewUsecaseImpl(supplierRegistry map[string]HotelSupplier, cache Cacher) *UsecaseImpl {
	return &UsecaseImpl{
		supplierRegistry: supplierRegistry,
		cache:            cache,
	}
}

func (u *UsecaseImpl) GetHotels(ctx context.Context, hotelIDs []string, destinationID int) ([]entity.Hotel, error) {
	// optimisation: we can use cache to store the results of the previous call
	// in this demo, we use the cache when user provides only hotelIDs
	// for each hotelID, we need to determine the cache key
	// we can use the hotelID as the key
	var cachedHotels []entity.Hotel
	remainingHotelIDs := make([]string, 0, len(hotelIDs))
	if len(hotelIDs) > 0 && destinationID < 0 {
		// if destinationID is not provided, we can use the hotelID as the cache key
		for _, hotelID := range hotelIDs {
			cachedHotel, err := u.getHotelFromCache(hotelIDs[0])
			if err != nil {
				// if there is error getting data from cache, we add the id to the remaining hotelIDs
				remainingHotelIDs = append(remainingHotelIDs, hotelID)
				continue
			}
			cachedHotels = append(cachedHotels, cachedHotel)
		}
	} else {
		// if not using cache, add all hotelIDs to the remaining hotelIDs
		remainingHotelIDs = hotelIDs
	}

	if len(remainingHotelIDs) == 0 && len(hotelIDs) > 0 {
		// if there are no remaining hotelIDs, we can return the list of hotels immediately
		return cachedHotels, nil
	}

	// concurrently fetch data from all suppliers
	p := pool.NewWithResults[[]entity.Hotel]()
	for _, supplier := range u.supplierRegistry {
		supplier := supplier // capture the loop variable
		p.Go(func() []entity.Hotel {
			logger := log.With().Str("supplier", supplier.GetName()).Logger()
			logger.Debug().Msgf("Fetching hotels from supplier %s", supplier.GetName())
			supplierHotels, err := supplier.FetchHotels(ctx, remainingHotelIDs, destinationID)
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

	// uniquely merge the data from all suppliers and return the final list
	mergedHotels := mergeHotelData(allHotels)

	// set the cache for the retrieved hotels
	for _, hotel := range mergedHotels {
		u.cache.Set(hotel.ID, hotel, time.Minute)
	}

	// concatenate the mergedHotels with the cachedHotels, if any
	mergedHotels = append(mergedHotels, cachedHotels...)

	return mergedHotels, nil
}

// cleanHotelData performs some basic cleaning on the hotel data before returning it to the caller.
func cleanHotelData(hotels []entity.Hotel) []entity.Hotel {
	for i, hotel := range hotels {
		// sanitise whitespace from hotel data
		hotels[i] = entity.TrimSpaceFromHotel(hotel)

		// for hotel amenities, each supplier may have a different naming style.
		// e.g. WiFi, Wifi, wifi, etc..
		// we need to normalise the names to a common format of all lowercase
		for i, amenity := range hotel.Amenities.General {
			hotel.Amenities.General[i] = strings.ToLower(amenity)
		}
		for i, amenity := range hotel.Amenities.Room {
			hotel.Amenities.Room[i] = strings.ToLower(amenity)
		}
	}

	return hotels
}

// mergeHotelData uniquely merges the hotel data from all suppliers to remove duplicates.
// This is done using hotel IDs as the key.
func mergeHotelData(hotels []entity.Hotel) []entity.Hotel {
	// create a map to hold aggregated hotel data
	// the key is the hotel ID, and the value is the hotel data
	aggregatedHotels := make(map[string]entity.Hotel)
	for _, hotel := range hotels {
		// check if the hotel already exists in the map
		if existingHotel, exists := aggregatedHotels[hotel.ID]; exists {
			// merge data if hotel already exists in the map.
			// merging rules are defined in the data model layer.
			aggregatedHotels[hotel.ID] = entity.MergeHotelData(existingHotel, hotel)
		} else {
			// add new hotel to the map if it doesn't exist
			aggregatedHotels[hotel.ID] = hotel
		}
	}

	// convert the map back to a slice
	finalHotels := make([]entity.Hotel, 0, len(aggregatedHotels))
	for _, hotel := range aggregatedHotels {
		finalHotels = append(finalHotels, hotel)
	}

	return finalHotels
}

// getHotelFromCache returns the hotel data from the cache if it exists.
// If the hotel data is not found in the cache, returns error.
func (u *UsecaseImpl) getHotelFromCache(hotelID string) (entity.Hotel, error) {
	// check if the hotelID is already in the cache
	data, found := u.cache.Get(hotelID)
	if !found {
		// if not found, we return an error
		return entity.Hotel{}, errors.New("hotel data not found in cache")
	}

	// if found, we can directly return the data
	cachedHotel, ok := data.(entity.Hotel)
	if !ok {
		// if the data is not a hotel, we return an error and log an error
		log.Error().Str("hotelID", hotelID).Msg("Cache data is not a hotel")
		return entity.Hotel{}, errors.New("cache data is not a hotel")
	}

	return cachedHotel, nil
}
