package helpers

import (
	"testing"
)

// SuggestService takes a string toComplete and returns a list of services that contain the provided string.
//
// The function uses a pre-defined list of services to find the suggestions.
//
// The function is case-insensitive and returns a list of strings in the same case as the original list.
func TestSuggestService(t *testing.T) {
	tests := []struct {
		name         string
		toComplete   string
		expected     []string
		expectedSize int
	}{
		{
			name:         "matching services found",
			toComplete:   "Softw",
			expected:     []string{"Software Engineering", "Software Development"},
			expectedSize: 2,
		},
		{
			name:         "no matching services found",
			toComplete:   "xyz",
			expected:     []string{},
			expectedSize: 0,
		},
		{
			name:         "matching services found with different case",
			toComplete:   "bloc",
			expected:     []string{"Blockchain Development"},
			expectedSize: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SuggestService(tt.toComplete)

			if len(got) != tt.expectedSize {
				t.Errorf("SuggestService(%s) returned %d suggestions, expected %d", tt.toComplete, len(got), tt.expectedSize)
			}

			for _, s := range tt.expected {
				if !containsString(got, s) {
					t.Errorf("SuggestService(%s) missing expected suggestion %s", tt.toComplete, s)
				}
			}
		})
	}
}

func containsString(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
