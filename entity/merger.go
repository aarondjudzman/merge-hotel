package entity

import (
	"strconv"
	"strings"
)

// MergeHotelData merges the data from the newHotel into the existingHotel.
// It uses some simple rules to merge the data that are stated within the function.
func MergeHotelData(existingHotel Hotel, newHotel Hotel) Hotel {
	// pick longer hotel name, if equal then use existing hotel name
	if len(newHotel.Name) > len(existingHotel.Name) {
		existingHotel.Name = newHotel.Name
	}

	// pick latitude and longitude with the highest accuracy. If equal, use existing hotel data
	if countDecimalPlaces(newHotel.Location.Latitude) > countDecimalPlaces(existingHotel.Location.Latitude) {
		existingHotel.Location.Latitude = newHotel.Location.Latitude
	}
	if countDecimalPlaces(newHotel.Location.Longitude) > countDecimalPlaces(existingHotel.Location.Longitude) {
		existingHotel.Location.Longitude = newHotel.Location.Longitude
	}

	// use longer address
	if len(newHotel.Location.Address) > len(existingHotel.Location.Address) {
		existingHotel.Location.Address = newHotel.Location.Address
	}
	if len(newHotel.Location.City) > len(existingHotel.Location.City) {
		existingHotel.Location.City = newHotel.Location.City
	}

	// use full country name instead of country code
	// possible improvement: if Country is a two-letter code, we can use mapping from ISO 3166-1 alpha-2 to country name
	if len(newHotel.Location.Country) > len(existingHotel.Location.Country) {
		existingHotel.Location.Country = newHotel.Location.Country
	}

	// concatenate hotel descriptions
	existingHotel.Description = existingHotel.Description + " " + newHotel.Description

	// for amenities, we follow two simple rules:
	// 		we remove duplicates from the each categories.
	// 		if there are also duplicates across both general and room amenities, we prioritise the general amenity
	existingHotel.Amenities.General = uniquelyMergeTwoLists(existingHotel.Amenities.General, newHotel.Amenities.General)
	existingHotel.Amenities.Room = uniquelyMergeTwoLists(existingHotel.Amenities.Room, newHotel.Amenities.Room)
	generalAmenitiesMap := make(map[string]bool)
	for _, amenity := range existingHotel.Amenities.General {
		generalAmenitiesMap[amenity] = true
	}
	existingHotel.Amenities.Room = uniqueRoomAmenities(newHotel.Amenities.Room, generalAmenitiesMap)

	// concatenate images and deduplicate
	existingHotel.Images.Amenities = mergeImages(existingHotel.Images.Amenities, newHotel.Images.Amenities)
	existingHotel.Images.Rooms = mergeImages(existingHotel.Images.Rooms, newHotel.Images.Rooms)
	existingHotel.Images.Site = mergeImages(existingHotel.Images.Site, newHotel.Images.Site)

	// concatenate booking conditions
	existingHotel.BookingConditions = append(existingHotel.BookingConditions, newHotel.BookingConditions...)

	return existingHotel
}

// countDecimalPlaces returns the number of decimal places in a float64 value.
func countDecimalPlaces(value float64) int {
	// Convert the float to a string.
	s := strconv.FormatFloat(value, 'f', -1, 64)
	// Find the index of the decimal point.
	i := strings.Index(s, ".")
	if i > -1 {
		// Return the number of characters after the decimal point.
		return len(s) - i - 1
	}
	// Return 0 if there is no decimal point.
	return 0
}

// uniquelyMergeTwoLists combines two lists of amenities, removing duplicates.
func uniquelyMergeTwoLists(list1, list2 []string) []string {
	amenitySet := make(map[string]struct{}) // using empty struct because it occupies no space, suitable when just checking if the key exists

	// add list1 to the map
	for _, amenity := range list1 {
		amenitySet[amenity] = struct{}{}
	}
	// add list2 to the map
	for _, amenity := range list2 {
		amenitySet[amenity] = struct{}{}
	}

	// convert the map to a slice
	mergedList := make([]string, 0, len(amenitySet))
	for amenity := range amenitySet {
		mergedList = append(mergedList, amenity)
	}

	return mergedList
}

// uniqueRoomAmenities removes duplicates from the roomAmenities list and returns a new list with only unique amenities.
func uniqueRoomAmenities(roomAmenities []string, generalSet map[string]bool) []string {
	var uniqueRoom []string
	for _, amenity := range roomAmenities {
		if !generalSet[amenity] {
			uniqueRoom = append(uniqueRoom, amenity)
		}
	}
	return uniqueRoom
}

// mergeImages deduplicates and merges two slices of Images.
func mergeImages(existing, incoming []Image) []Image {
	uniqueImages := make(map[string]Image) // Using a map to ensure uniqueness by image link.

	// add all existing images to the map.
	for _, img := range existing {
		uniqueImages[img.Link] = img
	}

	// try to add incoming images; skip duplicates.
	for _, img := range incoming {
		if existingImg, exists := uniqueImages[img.Link]; !exists {
			if existingImg.Description != img.Description {
				uniqueImages[img.Link] = img
			}
		}

	}

	// convert the map back to a slice.
	mergedImages := make([]Image, 0, len(uniqueImages))
	for _, img := range uniqueImages {
		mergedImages = append(mergedImages, img)
	}

	return mergedImages
}
