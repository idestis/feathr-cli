package helpers

import (
	"testing"
)

// TestContainsInSlice tests the ContainsInSlice function.
// The function should return true if a given string is present in a slice of strings,
// and false otherwise. The test cases cover both positive and negative cases.
func TestContainsInSlice(t *testing.T) {
	tests := []struct {
		name  string
		s     string
		slice []string
		want  bool
	}{
		{
			name:  "string in slice",
			s:     "foo",
			slice: []string{"foo", "bar", "baz"},
			want:  true,
		},
		{
			name:  "string not in slice",
			s:     "qux",
			slice: []string{"foo", "bar", "baz"},
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ContainsInSlice(tt.s, tt.slice)
			if got != tt.want {
				t.Errorf("ContainsInSlice(%q, %v) = %v, want %v", tt.s, tt.slice, got, tt.want)
			}
		})
	}
}

// TestIsValidEmail tests the IsValidEmail function.
// The function should return true if a given string is a valid email address,
// and false otherwise. The test cases cover both positive and negative cases,
// as well as edge cases and invalid email formats.
func TestIsValidEmail(t *testing.T) {
	testCases := []struct {
		email    string
		expected bool
	}{
		{"user@example.com", true},
		{"user123@example.com", true},
		{"user+123@example.com", true},
		{"user@example..com", false},
		{"user@.example.com", false},
		{"user@.com", false},
		{"@example.com", false},
		{"", false},
	}

	for _, tc := range testCases {
		result := IsValidEmail(tc.email)

		if result != tc.expected {
			t.Errorf("Unexpected result for email '%s': expected %v, got %v", tc.email, tc.expected, result)
		}
	}
}

// TestFindMaxInt tests the FindMaxInt function.
func TestFindMaxInt(t *testing.T) {
	testCases := []struct {
		nums []int
		want int
		err  bool
	}{
		{[]int{1, 2, 3}, 3, false},
		{[]int{-1, 0, 1}, 1, false},
		{[]int{}, 0, true},
	}

	for _, tc := range testCases {
		got, err := FindMaxInt(tc.nums)
		if (err != nil) != tc.err {
			t.Errorf("FindMaxInt(%v) error: %v, want error: %v", tc.nums, err, tc.err)
		}
		if got != tc.want {
			t.Errorf("FindMaxInt(%v) = %v, want %v", tc.nums, got, tc.want)
		}
	}
}
