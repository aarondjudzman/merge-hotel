package main

import (
	"ascenda-hotel/entity"
	"context"
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
	// TODO: implement the business logic of getting hotels data.

	return []entity.Hotel{}, nil
}
