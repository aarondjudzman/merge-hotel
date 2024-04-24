package supplier

import (
	"ascenda-hotel/entity"
	"testing"

	"github.com/efficientgo/core/testutil"
)

func TestConvertPaperfliesResponseToHotels(t *testing.T) {
	tests := []struct {
		name     string
		input    []PaperfliesResponse
		expected []entity.Hotel
	}{
		{
			name: "Single Hotel",
			input: []PaperfliesResponse{
				{
					HotelID:       "1",
					DestinationID: 101,
					HotelName:     "Paperflies Hotel",
					Location: PaperfliesLocation{
						Address: "123 Boulevard Rd",
						Country: "Wonderland",
					},
					Details: "A beautiful getaway.",
					Amenities: PaperfliesAmenities{
						General: []string{"Pool", "Free WiFi"},
						Room:    []string{"Air Conditioning", "Mini Bar"},
					},
					Images: PaperfliesImages{
						Rooms: []PaperfliesImage{{Link: "http://example.com/room.jpg", Caption: "View of the room"}},
						Site:  []PaperfliesImage{{Link: "http://example.com/site.jpg", Caption: "View of the hotel"}},
					},
					BookingConditions: []string{"No pets allowed", "Check-in at 3 PM"},
				},
			},
			expected: []entity.Hotel{
				{
					ID:            "1",
					DestinationID: 101,
					Name:          "Paperflies Hotel",
					Location: entity.Location{
						Address: "123 Boulevard Rd",
						Country: "Wonderland",
					},
					Description: "A beautiful getaway.",
					Amenities: entity.Amenities{
						General: []string{"Pool", "Free WiFi"},
						Room:    []string{"Air Conditioning", "Mini Bar"},
					},
					Images: entity.Images{
						Rooms:     []entity.Image{{Link: "http://example.com/room.jpg", Description: "View of the room"}},
						Site:      []entity.Image{{Link: "http://example.com/site.jpg", Description: "View of the hotel"}},
						Amenities: []entity.Image{},
					},
					BookingConditions: []string{"No pets allowed", "Check-in at 3 PM"},
				},
			},
		},
		{
			name:  "Multiple Hotels",
			input: []PaperfliesResponse{
				// Add another PaperfliesResponse similar to the first one but with different details
			},
			expected: []entity.Hotel{
				// Expected result for multiple hotels
			},
		},
		{
			name:     "Empty Input",
			input:    []PaperfliesResponse{},
			expected: []entity.Hotel{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convertPaperfliesResponseToHotels(tt.input)
			testutil.Equals(t, tt.expected, got)
		})
	}
}
