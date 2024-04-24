package supplier

import (
	"testing"

	"merge-hotel/entity"

	"github.com/efficientgo/core/testutil"
)

func TestConvertPatagoniaResponseToHotels(t *testing.T) {
	tests := []struct {
		name     string
		input    []PatagoniaResponse
		expected []entity.Hotel
	}{
		{
			name: "Basic Conversion",
			input: []PatagoniaResponse{
				{
					ID:          "1",
					Destination: 100,
					Name:        "Patagonia Hotel",
					Lat:         45.4215,
					Lng:         -75.6919,
					Address:     newString("123 Main St"),
					Info:        newString("A luxurious place."),
					Amenities:   []string{"Pool", "Spa"},
					Images: PatagoniaImages{
						Rooms: []PatagoniaImage{{URL: "http://example.com/room.jpg", Description: "Deluxe Room"}},
						Site:  []PatagoniaImage{{URL: "http://example.com/site.jpg", Description: "Front View"}},
					},
				},
			},
			expected: []entity.Hotel{
				{
					ID:            "1",
					DestinationID: 100,
					Name:          "Patagonia Hotel",
					Location: entity.Location{
						Latitude:  45.4215,
						Longitude: -75.6919,
						Address:   "123 Main St",
						City:      "",
						Country:   "",
					},
					Description: "A luxurious place.",
					Amenities: entity.Amenities{
						Room:    []string{"Pool", "Spa"},
						General: []string{},
					},
					Images: entity.Images{
						Rooms:     []entity.Image{{Link: "http://example.com/room.jpg", Description: "Deluxe Room"}},
						Site:      []entity.Image{{Link: "http://example.com/site.jpg", Description: "Front View"}},
						Amenities: []entity.Image{},
					},
					BookingConditions: []string{},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convertPatagoniaResponseToHotels(tt.input)
			testutil.Equals(t, tt.expected, got)
		})
	}
}
