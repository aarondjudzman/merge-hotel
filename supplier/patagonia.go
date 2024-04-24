package supplier

import (
	"context"
	"net/http"
	"time"

	"merge-hotel/entity"

	"github.com/carlmjohnson/requests"
)

// Patagonia is a supplier that fetches hotel data from Patagonia API.
type Patagonia struct {
	client  *http.Client
	address string
}

// NewPatagonia creates a new Patagonia supplier with the given endpoint address.
func NewPatagonia(address string) *Patagonia {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100

	return &Patagonia{
		client: &http.Client{
			Timeout: 2 * time.Second,
		},
		address: address,
	}
}

// PatagoniaResponse represents a hotel data response from Patagonia API.
type PatagoniaResponse struct {
	ID          string          `json:"id"`
	Destination int             `json:"destination"`
	Name        string          `json:"name"`
	Lat         float64         `json:"lat"`
	Lng         float64         `json:"lng"`
	Address     *string         `json:"address"`             // Using pointer to handle null
	Info        *string         `json:"info"`                // Using pointer to handle null
	Amenities   []string        `json:"amenities,omitempty"` // Using omitempty to handle empty slices
	Images      PatagoniaImages `json:"images"`
}

// PatagoniaImages defines the structure for images in the Patagonia API response.
type PatagoniaImages struct {
	Rooms     []PatagoniaImage `json:"rooms"`
	Site      []PatagoniaImage `json:"site"`
	Amenities []PatagoniaImage `json:"amenities"`
}

// PatagoniaImage defines the structure for a single image in the Patagonia API response.
type PatagoniaImage struct {
	URL         string `json:"url"`
	Description string `json:"description"`
}

// convertPatagoniaResponseToHotels converts PatagoniaResponse data to the common Hotel struct.
func convertPatagoniaResponseToHotels(responses []PatagoniaResponse) []entity.Hotel {
	hotels := make([]entity.Hotel, len(responses))
	for i, response := range responses {
		hotels[i] = entity.Hotel{
			ID:            response.ID,
			DestinationID: response.Destination,
			Name:          response.Name,
			Location:      convertPatagoniaLocations(response),
			Description:   derefString(response.Info), // handle possible null string
			Amenities:     convertPatagoniaAmenities(response.Amenities),
			Images: entity.Images{
				Rooms:     convertPatagoniaImages(response.Images.Rooms),
				Site:      convertPatagoniaImages(response.Images.Site),
				Amenities: convertPatagoniaImages(response.Images.Amenities),
			},
			BookingConditions: []string{}, // assuming no booking conditions data is directly available from Patagonia's response
		}
	}

	return hotels
}

// convertPatagoniaLocation converts PatagoniaResponse location data to the common Location struct.
func convertPatagoniaLocations(response PatagoniaResponse) entity.Location {
	return entity.Location{
		Latitude:  response.Lat,
		Longitude: response.Lng,
		Address:   derefString(response.Address), // handle possible null string
		City:      "",                            // assuming no city data is directly available from Patagonia's response
		Country:   "",                            // assuming no country data is directly available from Patagonia's response
	}
}

// convertPatagoniaAmenities converts PatagoniaResponse amenities data to Amenities in the common data model.
func convertPatagoniaAmenities(amenities []string) entity.Amenities {
	return entity.Amenities{
		Room:    amenities,  // assuming all patagonia amenities are considered as room amenities
		General: []string{}, // assuming no specific general amenities are listed in Patagonia's response
	}
}

// convertPatagoniaImages converts PatagoniaResponse images data to Images in the common data model.
func convertPatagoniaImages(images []PatagoniaImage) []entity.Image {
	convertedImages := make([]entity.Image, len(images))
	for i, image := range images {
		convertedImages[i] = entity.Image{
			Link:        image.URL,
			Description: image.Description,
		}
	}
	return convertedImages
}

func (p *Patagonia) FetchHotels(ctx context.Context, hotelIDs []string, destinationID int) ([]entity.Hotel, error) {
	// fetch hotels from the API
	var res []PatagoniaResponse
	err := requests.
		URL(p.address).
		ToJSON(&res).
		Client(p.client).
		Fetch(ctx)
	if err != nil {
		return []entity.Hotel{}, err
	}

	// filter hotels based on destination ID and hotel IDs if provided
	var filteredHotels []PatagoniaResponse
	hotelIDSet := make(map[string]bool) // create a map for quick lookup of HotelIDs
	if len(hotelIDs) > 0 {              // if hotelIDs is not empty, fill the map with the provided hotel IDs
		for _, id := range hotelIDs {
			hotelIDSet[id] = true
		}
	}

	for _, hotel := range res {
		// determine if the hotel should be included based on the provided parameters
		shouldIncludeDestination := (destinationID < 0) || hotel.Destination == destinationID
		shouldIncludeHotel := (len(hotelIDs) <= 0) || hotelIDSet[hotel.ID]

		if shouldIncludeDestination && shouldIncludeHotel {
			filteredHotels = append(filteredHotels, hotel)
		}
	}

	return convertPatagoniaResponseToHotels(filteredHotels), nil
}

// GetName returns the name of the supplier.
func (p *Patagonia) GetName() string {
	return "Patagonia"
}
