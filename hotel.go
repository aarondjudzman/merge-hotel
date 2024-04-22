package main

// Hotel represents the data model for a hotel.
// This data model is based on the response format for our API.
type Hotel struct {
	ID                string    `json:"id"`
	DestinationID     string    `json:"destination_id"`
	Name              string    `json:"name"`
	Location          Location  `json:"location"`
	Description       string    `json:"description"`
	Amenities         Amenities `json:"amenities"`
	Images            Images    `json:"images"`
	BookingConditions []string  `json:"booking_conditions"`
}

type Location struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
	Address   string  `json:"address"`
	City      string  `json:"city"`
	Country   string  `json:"country"`
}

type Amenities struct {
	General []string `json:"general"`
	Room    []string `json:"room"`
}

// Images are divided into categories
type Images struct {
	Rooms     []Image `json:"rooms"`
	Site      []Image `json:"site"`
	Amenities []Image `json:"amenities"`
}

// Image represents the attribute of an image.
// This is the same for every type of image
type Image struct {
	Link        string `json:"link"`
	Description string `json:"description"`
}
