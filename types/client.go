package types

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/idestis/feathr-cli/helpers"
	"github.com/spf13/cobra"
)

// Client represents a client of the business.
type Client struct {
	// ID is the unique identifier of the client.
	ID uint `survey:"id" yaml:"id"`

	// Name is the name of the client.
	Name string `survey:"name" json:"name"`

	// Address is the address of the client.
	Address string `survey:"address" json:"address"`

	// IBAN is a string that represents the user's International Bank Account Number (IBAN).
	IBAN string `survey:"iban" json:"iban"`

	// Bank is a string that represents the user's bank details.
	Bank string `survey:"bank" json:"bank"`

	// Email is the email address(es) of the client.
	Email []string `survey:"email" json:"email"`

	// Currency is the currency of the client.
	Currency string `survey:"currency" json:"currency"`
}

func (c *Client) WriteClientInfo(dataDir string, storageType string) error {
	ids, _ := GetClientIDs(dataDir, storageType)
	if len(ids) == 0 {
		c.ID = 1
	} else {
		id, _ := helpers.FindMaxInt(ids)
		c.ID = uint(id + 1)
	}
	switch storageType {
	case "file":
		err := os.MkdirAll(fmt.Sprintf("%s/data/clients/%d/invoices", dataDir, c.ID), 0700)
		if err != nil {
			return fmt.Errorf("failed to create client directory: %v", err)
		}

		filePath := fmt.Sprintf("%s/data/clients/%d/info.json", dataDir, c.ID)
		file, err := os.Create(filePath)
		if err != nil {
			return fmt.Errorf("failed to create file: %v", err)
		}
		defer file.Close()

		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(c); err != nil {
			return fmt.Errorf("failed to encode client info: %v", err)
		}

	case "sqlite":
		db, err := sql.Open("sqlite3", fmt.Sprintf("%s/feather-cli.sql", dataDir))
		if err != nil {
			return fmt.Errorf("failed to open SQLite database: %v", err)
		}
		defer db.Close()

		_, err = db.Exec("INSERT INTO clients (id, name, address, email) VALUES (?, ?, ?, ?, ?)", c.ID, c.Name, c.Address, c.Email)
		if err != nil {
			return fmt.Errorf("failed to insert client info into SQLite database: %v", err)
		}

	default:
		return fmt.Errorf("invalid storage type: %s", storageType)
	}

	return nil
}

// GetClientIDs returns a list of client IDs by scanning the directory
// where client data is stored.
func GetClientIDs(dataDir string, storageType string) ([]int, error) {
	var clientIDs []int

	switch storageType {
	case "file":
		clientsDir := filepath.Join(dataDir, "data/clients")
		files, err := ioutil.ReadDir(clientsDir)
		if err != nil {
			return nil, err
		}

		for _, f := range files {
			if f.IsDir() {
				id, _ := strconv.Atoi(f.Name())
				clientIDs = append(clientIDs, id)
			}
		}

	case "sqlite":
		db, err := sql.Open("sqlite3", fmt.Sprintf("%s/feather-cli.sql", dataDir))
		cobra.CheckErr(err)
		rows, err := db.Query("SELECT id FROM clients")
		cobra.CheckErr(err)
		defer rows.Close()

		for rows.Next() {
			var id int
			err := rows.Scan(&id)
			if err != nil {
				return nil, err
			}
			clientIDs = append(clientIDs, id)
		}

	default:
		return nil, fmt.Errorf("unknown storage type: %s", storageType)
	}

	return clientIDs, nil
}

// ReadClientInfo reads the client info for the given client ID from the configured storage.
// Returns an error if the client info cannot be read.
func ReadClientInfo(id int, dataDir string, storageType string) (Client, error) {
	var clientInfo Client

	switch storageType {
	case "file":
		// Build the path to the client info file.
		clientInfoFile := fmt.Sprintf("%s/data/clients/%d/info.json", dataDir, id)

		// Open the client info file.
		clientInfoBytes, err := ioutil.ReadFile(clientInfoFile)
		if err != nil {
			return clientInfo, err
		}

		// Unmarshal the client info data into a ClientInfo struct.
		err = json.Unmarshal(clientInfoBytes, &clientInfo)
		if err != nil {
			return clientInfo, err
		}

	case "sqlite":
		// Open a connection to the SQLite database.
		db, err := sql.Open("sqlite3", fmt.Sprintf("%s/feather-cli.sql", dataDir))
		if err != nil {
			return clientInfo, err
		}
		defer db.Close()

		// Prepare a query to select the client info for the given ID.
		query := "SELECT name, email, phone FROM clients WHERE id = ?"

		// Query the database and scan the result into a ClientInfo struct.
		err = db.QueryRow(query, id).Scan(&clientInfo.Name, &clientInfo.Email)
		if err != nil {
			return clientInfo, err
		}

	default:
		return clientInfo, fmt.Errorf("unknown storage type %q", storageType)
	}

	return clientInfo, nil
}
