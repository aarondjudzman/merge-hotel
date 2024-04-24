package entity

import (
	"strings"
)

// TrimSpaceFromHotel removes whitespace from suitable string fields in the hotel data.
func TrimSpaceFromHotel(hotel Hotel) Hotel {
	hotel.Name = strings.TrimSpace(hotel.Name)
	hotel.Location.Address = strings.TrimSpace(hotel.Location.Address)
	hotel.Location.City = strings.TrimSpace(hotel.Location.City)
	hotel.Location.Country = strings.TrimSpace(hotel.Location.Country)
	hotel.Description = strings.TrimSpace(hotel.Description)
	for i, amenity := range hotel.Amenities.General {
		hotel.Amenities.General[i] = strings.TrimSpace(amenity)
	}
	for i, amenity := range hotel.Amenities.Room {
		hotel.Amenities.Room[i] = strings.TrimSpace(amenity)
	}
	for i, image := range hotel.Images.Rooms {
		hotel.Images.Rooms[i].Description = strings.TrimSpace(image.Description)
	}
	for i, image := range hotel.Images.Site {
		hotel.Images.Site[i].Description = strings.TrimSpace(image.Description)
	}
	for i, image := range hotel.Images.Amenities {
		hotel.Images.Amenities[i].Description = strings.TrimSpace(image.Description)
	}
	for i, bookingCondition := range hotel.BookingConditions {
		hotel.BookingConditions[i] = strings.TrimSpace(bookingCondition)
	}

	return hotel
}
