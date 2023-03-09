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
	thousands "github.com/floscodes/golang-thousands"
	"github.com/idestis/feathr-cli/assets"
	"github.com/signintech/gopdf"
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
		unitPrice, _ := thousands.Separate(item.UnitPrice, "en")
		total, _ := thousands.Separate(item.UnitPrice*item.Quantity, "en")
		t.AddLine(item.Description, item.Quantity, unitPrice, total)
	}

	fmt.Printf("Invoice #%v \n\n", invoice.Number)
	fmt.Printf("Issued: %v\n", invoice.Issued.Format("2006-01-02"))
	fmt.Printf("Due: %v\n\n", invoice.Due.Format("2006-01-02"))
	if !invoice.Sent.IsZero() {
		fmt.Printf("Sent: %v\n\n", invoice.Sent.Format("2006-01-02"))
	}
	t.Print()
	if invoice.Notes != "" {
		fmt.Println()
		fmt.Println(invoice.Notes)
	}
}

func (invoice *Invoice) GeneratePDF(client Client, profile Profile, dataDir string) (error, string) {
	path := filepath.Join(dataDir, "generated")
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return fmt.Errorf("error creating generated folder: %v", err), ""
	}
	pdf := gopdf.GoPdf{}
	lineHeight := 16
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4}) //595.28, 841.89 = A4
	pdf.AddPage()
	data, err := assets.Asset("public/Roboto-Regular.ttf")
	if err != nil {
		return fmt.Errorf("error finding font: %v", err), ""
	}
	err = pdf.AddTTFFontData("Roboto", data)
	if err != nil {
		return fmt.Errorf("error finding font: %v", err), ""
	}
	err = pdf.SetFont("Roboto", "", 18)
	if err != nil {
		return fmt.Errorf("error setting font: %v", err), ""
	}
	pdf.AddPage()
	pdf.SetXY(25, 25)              // Starting possition is x=25, y=25
	pdf.SetTextColor(235, 100, 50) // Set Title Color
	pdf.Cell(nil, fmt.Sprintf("Invoice %v", invoice.Number))
	err = pdf.SetFont("Roboto", "", 13) // Reduce font size
	if err != nil {
		return fmt.Errorf("error reducing font size: %v", err), ""
	}
	pdf.SetXY(400, 25)
	pdf.SetTextColor(200, 200, 200)
	pdf.Cell(nil, "Issued:")
	pdf.SetXY(400, 40)
	pdf.Cell(nil, "Due:")
	pdf.SetTextColor(0, 0, 0) // Back to black color
	pdf.SetXY(450, 25)
	pdf.Cell(nil, invoice.Sent.Format("January 2, 2006"))
	pdf.SetXY(450, 40)
	pdf.Cell(nil, invoice.Due.Format("January 2, 2006"))
	// Client Information
	pdf.SetXY(30, 100)
	pdf.Cell(nil, client.Name)
	addressLines := strings.Split(client.Address, "\n")
	y := pdf.GetY()
	for _, line := range addressLines {
		pdf.SetX(30)
		pdf.SetY(y + float64(lineHeight))
		pdf.MultiCell(&gopdf.Rect{W: 250, H: float64(lineHeight)}, line)
		y += float64(lineHeight)
	}
	pdf.SetXY(300, 100)
	pdf.Cell(nil, profile.Name)
	pdf.SetXY(300, 115)
	profileAddressLines := strings.Split(profile.Address, "\n")
	y = pdf.GetY()
	for _, line := range profileAddressLines {
		pdf.SetX(300)
		pdf.SetY(y + float64(lineHeight))
		pdf.MultiCell(&gopdf.Rect{W: 250, H: float64(lineHeight)}, line)
		y += float64(lineHeight)
	}
	// pdf.MultiCell(&gopdf.Rect{W: 200, H: 40}, profile.Address)
	pdf.SetXY(300, pdf.GetY()+float64(lineHeight))
	pdf.Cell(nil, profile.IBAN)
	bankLines := strings.Split(profile.Bank, "\n")
	y = pdf.GetY() + float64(lineHeight)
	for _, line := range bankLines {
		pdf.SetX(300)
		pdf.SetY(y + float64(lineHeight))
		pdf.MultiCell(&gopdf.Rect{W: 250, H: float64(lineHeight)}, line)
		y += float64(lineHeight)
	}
	y = pdf.GetY() + float64(50)
	pdf.SetLineWidth(1)
	pdf.SetStrokeColor(200, 200, 200)
	qtyX := 300.0
	priceX := 400.0
	totalX := 500.0
	pdf.SetTextColor(200, 200, 200)
	pdf.SetXY(qtyX, y)
	pdf.Cell(nil, "Qty.")
	pdf.SetXY(priceX, y)
	pdf.Cell(nil, "Price")
	pdf.SetXY(totalX, y)
	pdf.Cell(nil, "Total")
	y += float64(lineHeight)
	pdf.Line(10, y, 585, y)
	pdf.SetTextColor(0, 0, 0)
	y = pdf.GetY() + float64(lineHeight) + 5
	for _, item := range invoice.Items {
		// TODO: Make this aproach to check if items are too long and we need a new page
		pdf.SetXY(15, y)
		pdf.MultiCell(&gopdf.Rect{W: 250, H: float64(lineHeight)}, item.Description)
		pdf.SetXY(qtyX, y)
		pdf.Cell(nil, fmt.Sprint(item.Quantity))
		pdf.SetXY(priceX, y)
		unitPrice, _ := thousands.Separate(item.UnitPrice, "en")
		pdf.Cell(nil, unitPrice)
		pdf.SetXY(totalX, y)
		itemTotal, _ := thousands.Separate(item.UnitPrice*item.Quantity, "en")
		pdf.Cell(nil, itemTotal)
		pdf.Line(10, y+float64(lineHeight), 585, y+float64(lineHeight))
		y += float64(lineHeight) + 5
	}
	y = pdf.GetY() + float64(lineHeight) + float64(lineHeight)
	pdf.SetXY(400, y)
	pdf.SetTextColor(200, 200, 200)
	pdf.Cell(nil, "Subtotal")
	pdf.SetXY(totalX, y)
	pdf.SetTextColor(235, 100, 50)
	subtotal, _ := thousands.Separate(invoice.GetInvoiceTotal(), "en")
	pdf.Cell(nil, subtotal)
	y = 670
	pdf.SetXY(15, y)
	pdf.SetTextColor(0, 0, 0)
	// Achieve word wrapping by splitting the text into multiple lines
	notes, _ := pdf.SplitTextWithWordWrap(invoice.Notes, 560)
	for _, line := range notes {
		pdf.SetXY(15, y)
		pdf.MultiCell(&gopdf.Rect{W: 560, H: float64(lineHeight)}, line)
		y += float64(lineHeight)
	}
	y = pdf.GetY() + (float64(lineHeight) * 2)
	pdf.Line(10, y, 585, y)
	y += (float64(lineHeight) * 2)
	pdf.SetXY(15, y)
	pdf.SetTextColor(200, 200, 200)
	pdf.Cell(nil, "Currency")
	pdf.SetXY(80, y)
	pdf.SetTextColor(0, 0, 0)
	pdf.Cell(nil, strings.ToUpper(client.Currency))
	// Write the invoice to a file
	invoicePath := filepath.Join(path, fmt.Sprintf("%v.pdf", invoice.ID))
	err = pdf.WritePdf(invoicePath) //TODO: make this dynamic
	if err != nil {
		return fmt.Errorf("error writing pdf: %v", err), ""
	}
	return nil, invoicePath
}

func (invoice *Invoice) Delete(dataDir string, storageType string) error {
	switch storageType {
	case "file":
		invoicePath := filepath.Join(dataDir, "data", "clients", fmt.Sprintf("%d", invoice.ClientID), "invoices", fmt.Sprintf("%v.json", invoice.ID))
		err := os.Remove(invoicePath)
		if err != nil {
			return fmt.Errorf("error deleting invoice: %v", err)
		}
		return nil
	case "sqlite":
		db, err := sql.Open("sqlite3", filepath.Join(dataDir, "feathr-cli.sql"))
		if err != nil {
			return fmt.Errorf("error opening database: %v", err)
		}
		defer db.Close()
		_, err = db.Query("DELETE FROM invoices WHERE id = ?", invoice.ID)
		if err != nil {
			return fmt.Errorf("error deleting invoice: %v", err)
		}
		return nil
	default:
		return fmt.Errorf("unknown storage type: %v", storageType)
	}
}
