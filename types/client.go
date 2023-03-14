package types

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/idestis/feathr-cli/helpers"
)

// Client represents a client of the business.
type Client struct {
	// ID is the unique identifier of the client.
	ID uint `survey:"id" json:"id"`

	// Name is the name of the client.
	Name string `survey:"name" json:"name"`

	// Address is the address of the client.
	Address string `survey:"address" json:"address"`

	// Bank is a string that represents the user's bank details.
	Bank string `survey:"bank" json:"bank"`

	// Email is the email address(es) of the client.
	Email []string `survey:"email" json:"email"`

	// Currency is the currency of the client.
	Currency string `survey:"currency" json:"currency"`
}

func (client *Client) WriteClientInfo(dataDir string) error {

	// During the filesystem as main storage, we have to search for last client ID first,
	// the next iteration will have caching to quickly fetch last ID instead of double work.
	ids, _ := GetClientIDs(dataDir)
	if len(ids) == 0 {
		client.ID = 1
	} else {
		id, _ := helpers.FindMaxInt(ids)
		client.ID = uint(id + 1)
	}
	err := os.MkdirAll(fmt.Sprintf("%s/data/clients/%d/invoices", dataDir, client.ID), 0700)
	if err != nil {
		return fmt.Errorf("failed to create client directory: %v", err)
	}

	filePath := fmt.Sprintf("%s/data/clients/%d/info.json", dataDir, client.ID)
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(client); err != nil {
		return fmt.Errorf("failed to encode client info: %v", err)
	}
	return nil
}

// GetClientIDs returns a list of client IDs by scanning the directory
// where client data is stored.
func GetClientIDs(dataDir string) ([]int, error) {
	var clientIDs []int

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
	return clientIDs, nil
}

// ReadClientInfo reads the client info for the given client ID from the storage.
// Returns an error if the client info cannot be read.
func ReadClientInfo(id int, dataDir string) (Client, error) {
	client := Client{
		ID: uint(id),
	}

	// Build the path to the client info file.
	clientInfoFile := fmt.Sprintf("%s/data/clients/%d/info.json", dataDir, client.ID)

	// Open the client info file.
	clientInfoBytes, err := ioutil.ReadFile(clientInfoFile)
	if err != nil {
		return client, err
	}

	// Unmarshal the client info data into a ClientInfo struct.
	err = json.Unmarshal(clientInfoBytes, &client)
	if err != nil {
		return client, err
	}

	return client, nil
}

func (client *Client) Print() error {
	fmt.Println("\033[32mAddress:\033[0m", client.Address)
	fmt.Println("\033[32mEmails:\033[0m", strings.Join(client.Email, ", "))
	if client.Bank != "" {
		fmt.Println("\033[32mBank Details:\033[0m", client.Bank)
	}
	return nil
}
