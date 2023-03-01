package types

import (
	"database/sql"
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

func (p *Profile) Exists(dataDir string, storageType string) bool {
	if storageType == "file" {
		filePath := filepath.Join(dataDir, "profile.json")
		if _, err := os.Stat(filePath); err == nil {
			return true
		}
	} else if storageType == "sqlite" {
		// check if the profile exists in the sqlite database
		// implementation details depend on your sqlite setup
	}
	return false
}

func (p *Profile) Load(dataDir string, storageType string) error {
	switch storageType {
	case "file":
		profilePath := filepath.Join(dataDir, "profile.json")
		if _, err := os.Stat(profilePath); os.IsNotExist(err) {
			// Profile file does not exist, return an error
			return errors.New("profile file does not exist")
		}
		data, err := ioutil.ReadFile(profilePath)
		if err != nil {
			return err
		}
		err = json.Unmarshal(data, p)
		if err != nil {
			return err
		}
	case "sqlite":
		db, err := sql.Open("sqlite3", filepath.Join(dataDir, "feather-cli.sqlite"))
		if err != nil {
			return err
		}
		defer db.Close()

		query := "SELECT id, name, currency, default_due, address, iban, bank_details FROM profile"
		rows, err := db.Query(query)
		if err != nil {
			return err
		}
		defer rows.Close()

		if !rows.Next() {
			// No rows returned, return an error
			return errors.New("profile not found in database")
		}

		err = rows.Scan(&p.Name, &p.Currency, &p.Due, &p.Address, &p.IBAN, &p.Bank)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid storage type: %s", storageType)
	}

	return nil
}

// Write writes the profile information to either a file or SQLite database.
func (p *Profile) Write(dataDir string, storageType string) error {
	if storageType == "file" {
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
	} else if storageType == "sqlite" {
		// TODO: Write the profile information to SQLite.
	} else {
		return fmt.Errorf("unknown storage type: %s", storageType)
	}

	return nil
}
