package supplier

import "testing"

func TestDerefFloat64(t *testing.T) {
	tests := []struct {
		input    *float64
		name     string
		expected float64
	}{
		{
			name:     "Non-nil float64",
			input:    newFloat64(3.14),
			expected: 3.14,
		},
		{
			name:     "Nil float64",
			input:    nil,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := derefFloat64(tt.input); got != tt.expected {
				t.Errorf("derefFloat64() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestDerefString(t *testing.T) {
	tests := []struct {
		name     string
		input    *string
		expected string
	}{
		{
			name:     "Non-nil string",
			input:    newString("hello"),
			expected: "hello",
		},
		{
			name:     "Nil string",
			input:    nil,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := derefString(tt.input); got != tt.expected {
				t.Errorf("derefString() = %v, want %v", got, tt.expected)
			}
		})
	}
}
func newString(val string) *string {
	return &val
}
