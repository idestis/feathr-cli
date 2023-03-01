package helpers

import (
	"errors"
	"net/mail"
)

// ContainsInSlice checks if a string is present in a slice of strings.
func ContainsInSlice(s string, slice []string) bool {
	// Convert the slice to a set to make lookups more efficient.
	set := make(map[string]struct{})
	for _, v := range slice {
		set[v] = struct{}{}
	}

	// Check if the string is in the set.
	_, ok := set[s]
	return ok
}

// ContainsIntInSlice checks if a int is present in a slice of ints.
func ContainsIntInSlice(s int, slice []int) bool {
	// Convert the slice to a set to make lookups more efficient.
	set := make(map[int]struct{})
	for _, v := range slice {
		set[v] = struct{}{}
	}

	// Check if the string is in the set.
	_, ok := set[s]
	return ok
}

// IsValidEmail returns true if the given string is a valid email address, and false otherwise.
// The function uses a built-in mail package to perform the validation.
func IsValidEmail(email string) bool {
	// Parse the email address using the mail.ParseAddress function.
	parsed, err := mail.ParseAddress(email)

	// If there was an error parsing the email address, it is not valid.
	if err != nil {
		return false
	}

	// If the parsed address is the same as the input address, it is valid.
	return parsed.Address == email
}

// FindMaxInt returns the maximum integer in a slice of integers.
func FindMaxInt(nums []int) (int, error) {
	if len(nums) == 0 {
		return 0, errors.New("cannot find maximum of empty slice")
	}

	max := nums[0]
	for _, num := range nums[1:] {
		if num > max {
			max = num
		}
	}

	return max, nil
}
