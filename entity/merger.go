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

	// for amenities, we follow three simple rules:
	// 		- we remove duplicates from the each categories.
	// 	    - some supplier may also return amenities in the form of concatenated strings
	// 		  e.g. "DryCleaning, BathTub", while others return as "dry cleaning, bath tub"
	// 		  we will check for duplicates such as "dry cleaning" and "drycleaning" and remove the duplicates,
	// 		  prioritising the one with space in between
	// 		- if there are also duplicates across both general and room amenities, we prioritise the room amenity

	existingHotel.Amenities.General = uniquelyMergeTwoLists(existingHotel.Amenities.General, newHotel.Amenities.General)
	existingHotel.Amenities.Room = uniquelyMergeTwoLists(existingHotel.Amenities.Room, newHotel.Amenities.Room)
	existingHotel.Amenities.General = removeRoomAmenities(existingHotel.Amenities.Room, existingHotel.Amenities.General)

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
	amenitySet := make(map[string]string) // Map to hold the normalized string as key and preferred format as value

	// function to add items to the map
	addAmenities := func(list []string) {
		for _, amenity := range list {
			normalized := strings.ToLower(strings.ReplaceAll(amenity, " ", "")) // normalize the string by removing spaces and converting to lowercase
			if existing, exists := amenitySet[normalized]; exists {
				// check if the current entry has spaces and the existing one does not, replace it
				if strings.Contains(amenity, " ") && !strings.Contains(existing, " ") {
					amenitySet[normalized] = amenity
				}
			} else {
				// if it doesn't exist, add the new item
				amenitySet[normalized] = amenity
			}
		}
	}

	// add both lists to the map
	addAmenities(list1)
	addAmenities(list2)

	// convert the map to a slice
	mergedList := make([]string, 0, len(amenitySet))
	for _, amenity := range amenitySet {
		mergedList = append(mergedList, amenity)
	}

	return mergedList
}

// removeRoomAmenities checks if general amenities already exists in the room amenities list and returns a new list with only unique amenities.
func removeRoomAmenities(roomAmenities, generalAmenities []string) []string {
	roomAmenitiesMap := make(map[string]bool)
	for _, amenity := range roomAmenities {
		roomAmenitiesMap[amenity] = true
	}

	var uniqueGeneral []string
	for _, amenity := range generalAmenities {
		if !roomAmenitiesMap[amenity] {
			uniqueGeneral = append(uniqueGeneral, amenity)
		}
	}
	return uniqueGeneral
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
