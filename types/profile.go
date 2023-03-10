package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Profile represents a user's profile information.
type Profile struct {
	// Name is a string that represents the user's name.
	Name string `survey:"name" yaml:"name" json:"name"`

	// Currency is a string that represents the user's preferred currency.
	Currency string `survey:"currency" yaml:"currency" json:"currency"`

	// Due is a string that represents the user's preferred due date.
	Due string `survey:"default_due" yaml:"default_due" json:"default_due"`

	// Address is a string that represents the user's address.
	Address string `survey:"address" yaml:"address" json:"address"`

	// IBAN is a string that represents the user's International Bank Account Number (IBAN).
	IBAN string `survey:"iban" yaml:"iban" json:"iban"`

	// Bank is a string that represents the user's bank details.
	Bank string `survey:"bank_details" yaml:"bank_details" json:"bank_details"`
}

// Exists checks if the profile exists
func (p *Profile) Exists(dataDir string) bool {
	filePath := filepath.Join(dataDir, "profile.json")
	if _, err := os.Stat(filePath); err == nil {
		return true
	}
	return false
}

func (p *Profile) Load(dataDir string) error {
	profilePath := filepath.Join(dataDir, "profile.json")
	if _, err := os.Stat(profilePath); os.IsNotExist(err) {
		// Profile file does not exist, return an error
		return errors.New("profile file does not exist. have you run `feather-cli profile`?")
	}
	data, err := ioutil.ReadFile(profilePath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, p)
	if err != nil {
		return err
	}
	return nil
}

// Write writes the profile information to a file
func (p *Profile) Write(dataDir string) error {
	// Marshal the profile information to JSON.
	bytes, err := json.MarshalIndent(p, "", "\t")
	if err != nil {
		return fmt.Errorf("error marshaling profile: %w", err)
	}

	// Write the JSON to a file.
	filename := filepath.Join(dataDir, "profile.json")
	err = ioutil.WriteFile(filename, bytes, 0644)
	if err != nil {
		return fmt.Errorf("error writing profile file: %w", err)
	}
	return nil
}
