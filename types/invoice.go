package types

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/cheynewallace/tabby"
)

// Invoice represents an invoice for a client.
type Invoice struct {
	// ID is the unique identifier of the invoice.
	ID uint `json:"id"`

	// ClientID is the unique identifier of the client associated with the invoice.
	ClientID uint `json:"client_id"`

	// Number is the invoice number.
	Number string `survey:"number" json:"number"`

	// Issued is the date the invoice was issued.
	Issued time.Time `json:"issued"`

	// Due is the due date for the invoice.
	Due time.Time `survey:"due" json:"due"`

	// Sent is the date when the invoice was sent.
	Sent time.Time `json:"sent"`

	// Notes is any additional notes associated with the invoice.
	Notes string `survey:"notes" json:"notes"`

	// Items is a list of line items on the invoice.
	Items []Item `survey:"items" json:"items"`

	// Subtotal is subtotal of the invoice.
	Subtotal float32 `json:"subtotal"`
}

// Item represents a line item on an invoice.
type Item struct {
	// ID is the unique identifier of the item.
	ID uint `json:"id"`

	// InvoiceID is the unique identifier of the invoice the item belongs to.
	InvoiceID uint `json:"invoice_id"`

	// Description is a brief description of the item.
	Description string `survey:"description" json:"description"`

	// Quantity is the quantity of the item.
	Quantity float64 `survey:"quantity" json:"quantity"`

	// UnitPrice is the unit price of the item.
	UnitPrice float64 `survey:"unit_price" json:"unit_price"`
}

func GetInvoicesID(clientIDs []int, dataDir string, storageType string) (map[int]string, []int, error) {
	// Initialize variables to hold the invoice IDs and paths.
	invoiceIDs := make([]int, 0)
	invoicePaths := make(map[int]string)

	// Iterate over the client IDs.
	for _, id := range clientIDs {
		// Check if the storage type is "file".
		if storageType == "file" {
			// Construct the path to the invoices directory for this client.
			invoicesPath := filepath.Join(dataDir, "data", "clients", strconv.Itoa(id), "invoices")
			// Walk over the invoices directory to find invoice files.
			err := filepath.Walk(invoicesPath, func(path string, info os.FileInfo, err error) error {
				// Check if the file is a JSON file.
				if strings.HasSuffix(info.Name(), ".json") {
					// Extract the invoice ID from the file name.
					parts := strings.Split(info.Name(), ".")
					if len(parts) != 2 {
						return fmt.Errorf("invalid invoice file name: %s", info.Name())
					}
					invoiceID, err := strconv.Atoi(parts[0])
					if err != nil {
						return fmt.Errorf("invalid invoice file name: %s", info.Name())
					}

					// Add the invoice ID and path to the maps.
					invoiceIDs = append(invoiceIDs, invoiceID)
					invoicePaths[invoiceID] = strings.Replace(path, dataDir, "", 1)
				}
				return nil
			})
			if err != nil {
				return nil, nil, err
			}
		} else if storageType == "sqlite" {
			// Query the database to get the invoice IDs.
			db, err := sql.Open("sqlite3", filepath.Join(dataDir, "feathr-cli.sql"))
			if err != nil {
				return nil, nil, err
			}
			defer db.Close()
			rows, err := db.Query("SELECT id FROM invoices WHERE client_id = ?", id)
			if err != nil {
				return nil, nil, err
			}
			defer rows.Close()

			// Add the invoice IDs to the slice.
			for rows.Next() {
				var invoiceID int
				err := rows.Scan(&invoiceID)
				if err != nil {
					return nil, nil, err
				}
				invoiceIDs = append(invoiceIDs, invoiceID)
			}
			err = rows.Err()
			if err != nil {
				return nil, nil, err
			}
		} else {
			return nil, nil, fmt.Errorf("invalid storage type: %s", storageType)
		}
	}

	// Sort the slice of invoice IDs.
	sort.Ints(invoiceIDs)

	return invoicePaths, invoiceIDs, nil
}

func (invoice *Invoice) WriteInvoice(dataDir string, storageType string) error {
	if storageType == "file" {
		// Create the invoices directory for this client if it doesn't exist.
		invoicesDir := filepath.Join(dataDir, "data", "clients", fmt.Sprintf("%d", invoice.ClientID), "invoices")
		err := os.MkdirAll(invoicesDir, 0755)
		if err != nil {
			return err
		}

		// Marshal the invoice to JSON.
		data, err := json.Marshal(invoice)
		if err != nil {
			return err
		}

		// Write the JSON data to a file.
		invoicePath := filepath.Join(invoicesDir, fmt.Sprintf("%d.json", invoice.ID))
		err = ioutil.WriteFile(invoicePath, data, 0644)
		if err != nil {
			return err
		}

		return nil
	} else if storageType == "sqlite" {
		// Insert the invoice into the database.
		db, err := sql.Open("sqlite3", filepath.Join(dataDir, "feathr-cli.sql"))
		if err != nil {
			return err
		}
		defer db.Close()

		stmt, err := db.Prepare("INSERT INTO invoices (id, client_id, amount) VALUES (?, ?, ?)")
		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(invoice.ID, invoice.ClientID, invoice.Items)
		if err != nil {
			return err
		}

		return nil
	} else {
		return fmt.Errorf("invalid storage type: %s", storageType)
	}
}

func ReadInvoice(id int, path string, dataDir string, storageType string) (Invoice, error) {
	invoice := Invoice{}

	if storageType == "file" {
		// Read the JSON data from the file.
		invoicePath := filepath.Join(dataDir, path)
		data, err := ioutil.ReadFile(invoicePath)
		if err != nil {
			return invoice, err
		}

		// Unmarshal the JSON data into the invoice.
		err = json.Unmarshal(data, &invoice)
		if err != nil {
			return invoice, err
		}

		return invoice, nil
	} else if storageType == "sqlite" {
		// Query the database to get the invoice.
		db, err := sql.Open("sqlite3", filepath.Join(dataDir, "feathr-cli.sql"))
		if err != nil {
			return invoice, err
		}
		defer db.Close()

		row := db.QueryRow("SELECT id, client_id, number, issued, due, sent, notes FROM invoices WHERE id = ?", id)
		err = row.Scan(&invoice.ID, &invoice.ClientID, &invoice.Number, &invoice.Issued, &invoice.Due, &invoice.Sent, &invoice.Notes)
		if err != nil {
			return invoice, err
		}

		return invoice, nil
	} else {
		return invoice, fmt.Errorf("invalid storage type: %s", storageType)
	}
}

func (invoice *Invoice) GetInvoiceTotal() float64 {
	total := 0.0
	for _, item := range invoice.Items {
		total += item.UnitPrice * float64(item.Quantity)
	}
	return total
}

func (invoice *Invoice) Print() {
	t := tabby.New()
	t.AddHeader("Description", "Qty.", "Price", "Total")
	for _, item := range invoice.Items {
		t.AddLine(item.Description, item.Quantity, item.UnitPrice, item.UnitPrice*item.Quantity)
	}

	fmt.Printf("Invoice #%v \n\n", invoice.Number)
	fmt.Printf("Issued: %v\n", invoice.Issued.Format("2006-01-02"))
	fmt.Printf("Due: %v\n\n", invoice.Due.Format("2006-01-02"))
	if !invoice.Sent.IsZero() {
		fmt.Printf("Sent: %v\n\n", invoice.Sent.Format("2006-01-02"))
	}
	t.Print()
	if invoice.Notes != "" {
		fmt.Printf("\nNotes: %v\n", invoice.Notes)
	}
}
