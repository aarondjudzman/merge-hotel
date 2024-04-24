package supplier

import (
	"context"
	"net/http"
	"time"

	"merge-hotel/entity"

	"github.com/carlmjohnson/requests"
)

// Paperflies is a supplier that fetches hotel data from Paperflies API.
type Paperflies struct {
	client  *http.Client
	address string
}

// NewPaperflies creates a new Paperflies supplier with the given endpoint address.
func NewPaperflies(address string) *Paperflies {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100

	return &Paperflies{
		client: &http.Client{
			Timeout: 2 * time.Second,
		},
		address: address,
	}
}

// PaperfliesResponse represents a hotel data response from Paperflies API.
type PaperfliesResponse struct {
	HotelID           string              `json:"hotel_id"`
	DestinationID     int                 `json:"destination_id"`
	HotelName         string              `json:"hotel_name"`
	Location          PaperfliesLocation  `json:"location"`
	Details           string              `json:"details"`
	Amenities         PaperfliesAmenities `json:"amenities"`
	Images            PaperfliesImages    `json:"images"`
	BookingConditions []string            `json:"booking_conditions"`
}

// PaperfliesLocation defines the structure for location details in the Paperflies API response.
type PaperfliesLocation struct {
	Address string `json:"address"`
	Country string `json:"country"`
}

// PaperfliesAmenities defines the structure for amenities, split into general and room categories.
type PaperfliesAmenities struct {
	General []string `json:"general"`
	Room    []string `json:"room"`
}

// PaperfliesImages defines the structure for images in the Paperflies API response, categorized into rooms and site.
type PaperfliesImages struct {
	Rooms []PaperfliesImage `json:"rooms"`
	Site  []PaperfliesImage `json:"site"`
}

// PaperfliesImage defines the structure for a single image in the Paperflies API response.
type PaperfliesImage struct {
	Link    string `json:"link"`
	Caption string `json:"caption"`
}

// convertPaperfliesResponseToHotels converts PaperfliesResponse data to the common Hotel struct.
func convertPaperfliesResponseToHotels(responses []PaperfliesResponse) []entity.Hotel {
	hotels := make([]entity.Hotel, len(responses))
	for i, response := range responses {
		hotels[i] = entity.Hotel{
			ID:            response.HotelID,
			DestinationID: response.DestinationID,
			Name:          response.HotelName,
			Location:      convertPaperfliesLocations(response),
			Description:   response.Details,
			Amenities:     convertPaperfliesAmenities(response.Amenities),
			Images: entity.Images{
				Rooms:     convertPaperfliesImages(response.Images.Rooms),
				Site:      convertPaperfliesImages(response.Images.Site),
				Amenities: []entity.Image{}, // assuming no specific amenity images are listed in Paperflies's response
			},
			BookingConditions: response.BookingConditions,
		}
	}

	return hotels
}

// convertPaperfliesLocation converts PaperfliesResponse location data to the common Location struct.
func convertPaperfliesLocations(response PaperfliesResponse) entity.Location {
	return entity.Location{
		Address: response.Location.Address,
		Country: response.Location.Country,
	}
}

// convertPaperfliesAmenities converts PaperfliesResponse amenities data to Amenities in the common data model.
func convertPaperfliesAmenities(amenities PaperfliesAmenities) entity.Amenities {
	return entity.Amenities{
		General: amenities.General,
		Room:    amenities.Room,
	}
}

// convertPaperfliesImages converts PaperfliesResponse images data to Images in the common data model.
func convertPaperfliesImages(images []PaperfliesImage) []entity.Image {
	convertedImages := make([]entity.Image, len(images))
	for i, image := range images {
		convertedImages[i] = entity.Image{
			Link:        image.Link,
			Description: image.Caption,
		}
	}
	return convertedImages
}

// FetchHotels fetches hotels from the Paperflies API.
func (p *Paperflies) FetchHotels(ctx context.Context, hotelIDs []string, destinationID int) ([]entity.Hotel, error) {
	// fetch hotels from the API
	var res []PaperfliesResponse
	err := requests.
		URL(p.address).
		ToJSON(&res).
		Client(p.client).
		Fetch(ctx)
	if err != nil {
		return []entity.Hotel{}, err
	}

	// filter hotels based on destination ID and hotel IDs if provided
	var filteredHotels []PaperfliesResponse
	hotelIDSet := make(map[string]bool) // create a map for quick lookup of HotelIDs
	if len(hotelIDs) > 0 {              // if hotelIDs is not empty, fill the map with the provided hotel IDs
		for _, id := range hotelIDs {
			hotelIDSet[id] = true
		}
	}

	for _, hotel := range res {
		// determine if the hotel should be included based on the provided parameters
		shouldIncludeDestination := (destinationID < 0) || hotel.DestinationID == destinationID
		shouldIncludeHotel := (len(hotelIDs) <= 0) || hotelIDSet[hotel.HotelID]

		if shouldIncludeDestination && shouldIncludeHotel {
			filteredHotels = append(filteredHotels, hotel)
		}
	}

	return convertPaperfliesResponseToHotels(filteredHotels), nil
}

func (p *Paperflies) GetName() string {
	return "Paperflies"
}
