package supplier

import (
	"ascenda-hotel/entity"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/carlmjohnson/requests"
)

// Acme is a supplier that fetches hotel data from the Acme API.
type Acme struct {
	client  *http.Client
	address string
}

// NewAcme creates a new Acme supplier with the given endpoint address.
func NewAcme(address string) *Acme {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100

	return &Acme{
		client: &http.Client{
			Timeout: 2 * time.Second,
		},
		address: address,
	}
}

// AcmeResponse represents the response format for a hotel from the Acme API.
type AcmeResponse struct {
	Latitude      *float64 `json:"Latitude"`  // Using pointer to handle null
	Longitude     *float64 `json:"Longitude"` // Using pointer to handle null
	ID            string   `json:"Id"`
	Name          string   `json:"Name"`
	Address       string   `json:"Address"`
	City          string   `json:"City"`
	Country       string   `json:"Country"`
	PostalCode    string   `json:"PostalCode"`
	Description   string   `json:"Description"`
	Facilities    []string `json:"Facilities"`
	DestinationID int      `json:"DestinationId"`
}

// UnmarshalJSON helps in custom unmarshaling to handle potential issues with data types.
func (a *AcmeResponse) UnmarshalJSON(data []byte) error {
	type Alias AcmeResponse
	aux := &struct {
		*Alias
		Latitude  json.RawMessage `json:"Latitude"`
		Longitude json.RawMessage `json:"Longitude"`
	}{
		Alias: (*Alias)(a),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// handle Latitude
	if err := parseOptionalFloat(aux.Latitude, &a.Latitude); err != nil {
		return err
	}

	// handle Longitude
	if err := parseOptionalFloat(aux.Longitude, &a.Longitude); err != nil {
		return err
	}

	return nil
}

// parseOptionalFloat parses a JSON raw message that may be a number, null, string with quotes.
func parseOptionalFloat(raw json.RawMessage, target **float64) error {
	if len(raw) == 0 || string(raw) == "null" {
		*target = nil
		return nil
	}
	// check if the raw message is a number or a string with quotes
	var num float64
	if err := json.Unmarshal(raw, &num); err == nil {
		*target = &num
		return nil
	}
	var str string
	if err := json.Unmarshal(raw, &str); err == nil {
		if str == "" {
			*target = nil
			return nil
		}
		num, err = strconv.ParseFloat(str, 64)
		if err != nil {
			return errors.New("invalid number literal")
		}
		*target = &num
		return nil
	}
	return errors.New("invalid JSON for latitude or longitude")
}

// convertAcmeResponseToHotels converts AcmeResponse data to the common Hotel struct.
func convertAcmeResponseToHotels(responses []AcmeResponse) []entity.Hotel {
	hotels := make([]entity.Hotel, len(responses))
	for i, response := range responses {
		hotels[i] = entity.Hotel{
			ID:            response.ID,
			DestinationID: response.DestinationID,
			Name:          response.Name,
			Location:      convertAcmeLocations(response),
			Description:   response.Description,
			Amenities:     convertAcmeAmenities(response.Facilities),
			Images: entity.Images{
				Rooms:     []entity.Image{},
				Site:      []entity.Image{},
				Amenities: []entity.Image{},
			}, // assuming no image data is directly available from Acme's response
			BookingConditions: []string{}, // Assuming no booking conditions data is directly available from Acme's response
		}
	}

	return hotels
}

// convertAcmeLocation converts AcmeResponse location data to the common Location struct.
func convertAcmeLocations(response AcmeResponse) entity.Location {
	return entity.Location{
		// for both latitude and longitude, in the case of null or empty string, we use 0 as the value
		Latitude:  derefFloat64(response.Latitude), // convert *float64 to float64
		Longitude: derefFloat64(response.Longitude),
		Address:   response.Address,
		City:      response.City,
		Country:   response.Country,
	}
}

// convertAcmeAmenities converts AcmeResponse facilities data to Amenities in the common data model.
func convertAcmeAmenities(facilities []string) entity.Amenities {
	// acme amenities are not categorised, so it is assumed that all facilities are considered as general amenities
	return entity.Amenities{
		General: facilities,
		Room:    []string{},
	}
}

// FetchHotels fetches hotels from the Acme API.
func (a *Acme) FetchHotels(ctx context.Context, hotelIDs []string, destinationID int) ([]entity.Hotel, error) {
	// fetch hotels from the API
	var res []AcmeResponse
	err := requests.
		URL(a.address).
		ToJSON(&res).
		Client(a.client).
		Fetch(ctx)
	if err != nil {
		return []entity.Hotel{}, err
	}

	// filter hotels based on destination ID and hotel IDs if provided
	var filteredHotels []AcmeResponse
	hotelIDSet := make(map[string]bool) // create a map for quick lookup of HotelIDs
	if len(hotelIDs) > 0 {              // if hotelIDs is not empty, fill the map with the provided hotel IDs
		for _, id := range hotelIDs {
			hotelIDSet[id] = true
		}
	}

	for _, hotel := range res {
		// determine if the hotel should be included based on the provided parameters
		shouldIncludeDestination := (destinationID < 0) || hotel.DestinationID == destinationID
		shouldIncludeHotel := (len(hotelIDs) <= 0) || hotelIDSet[hotel.ID]

		if shouldIncludeDestination && shouldIncludeHotel {
			filteredHotels = append(filteredHotels, hotel)
		}
	}

	return convertAcmeResponseToHotels(filteredHotels), nil
}

// GetName returns the name of the supplier.
func (a *Acme) GetName() string {
	return "Acme"
}
