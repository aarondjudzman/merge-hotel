package supplier

import (
	"encoding/json"
	"testing"

	"ascenda-hotel/entity"

	"github.com/efficientgo/core/testutil"
)

func TestParseOptionalFloat(t *testing.T) {
	validNumber := json.RawMessage(`123.45`)
	var num *float64
	if err := parseOptionalFloat(validNumber, &num); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if num == nil || *num != 123.45 {
		t.Fatalf("Expected 123.45, got %v", num)
	}

	invalidNumber := json.RawMessage(`"abc"`)
	if err := parseOptionalFloat(invalidNumber, &num); err == nil {
		t.Fatalf("Expected an error for invalid number, got nil")
	}
}

func TestConvertAcmeResponseToHotels(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		input    []AcmeResponse
		expected []entity.Hotel
	}{
		{
			name: "Single Hotel",
			input: []AcmeResponse{
				{
					ID:            "1",
					DestinationID: 101,
					Name:          "Acme Hotel",
					Latitude:      newFloat64(40.7128),
					Longitude:     newFloat64(-74.0060),
					Address:       "123 Example St",
					City:          "New York",
					Country:       "USA",
					Description:   "A nice place.",
					Facilities:    []string{"Free WiFi", "Parking"},
				},
			},
			expected: []entity.Hotel{
				{
					ID:            "1",
					DestinationID: 101,
					Name:          "Acme Hotel",
					Location:      entity.Location{Latitude: 40.7128, Longitude: -74.0060, Address: "123 Example St", City: "New York", Country: "USA"},
					Description:   "A nice place.",
					Amenities:     entity.Amenities{General: []string{"Free WiFi", "Parking"}, Room: []string{}},
					Images: entity.Images{
						Rooms:     []entity.Image{},
						Site:      []entity.Image{},
						Amenities: []entity.Image{},
					},
					BookingConditions: []string{},
				},
			},
		},
		{
			name: "Multiple Hotels",
			input: []AcmeResponse{
				{
					ID:            "1",
					DestinationID: 102,
					Name:          "Acme Hotel One",
					Latitude:      newFloat64(34.0522),
					Longitude:     newFloat64(-118.2437),
					Address:       "456 Example Rd",
					City:          "Los Angeles",
					Country:       "USA",
					Description:   "Luxurious stay.",
					Facilities:    []string{"Pool", "Spa"},
				},
				{
					ID:            "2",
					DestinationID: 103,
					Name:          "Acme Hotel Two",
					Latitude:      newFloat64(37.7749),
					Longitude:     newFloat64(-122.4194),
					Address:       "789 Example Blvd",
					City:          "San Francisco",
					Country:       "USA",
					Description:   "Comfortable and central.",
					Facilities:    []string{"Free WiFi", "Gym"},
				},
			},
			expected: []entity.Hotel{
				{
					ID:            "1",
					DestinationID: 102,
					Name:          "Acme Hotel One",
					Location:      entity.Location{Latitude: 34.0522, Longitude: -118.2437, Address: "456 Example Rd", City: "Los Angeles", Country: "USA"},
					Description:   "Luxurious stay.",
					Amenities:     entity.Amenities{General: []string{"Pool", "Spa"}, Room: []string{}},
					Images: entity.Images{
						Rooms:     []entity.Image{},
						Site:      []entity.Image{},
						Amenities: []entity.Image{},
					},
					BookingConditions: []string{},
				},
				{
					ID:            "2",
					DestinationID: 103,
					Name:          "Acme Hotel Two",
					Location:      entity.Location{Latitude: 37.7749, Longitude: -122.4194, Address: "789 Example Blvd", City: "San Francisco", Country: "USA"},
					Description:   "Comfortable and central.",
					Amenities:     entity.Amenities{General: []string{"Free WiFi", "Gym"}, Room: []string{}},
					Images: entity.Images{
						Rooms:     []entity.Image{},
						Site:      []entity.Image{},
						Amenities: []entity.Image{},
					},
					BookingConditions: []string{},
				},
			},
		},
		{
			name:     "Empty Input",
			input:    []AcmeResponse{},
			expected: []entity.Hotel{},
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertAcmeResponseToHotels(tt.input)
			testutil.Equals(t, tt.expected, result)
		})
	}
}

func newFloat64(val float64) *float64 {
	return &val
}
