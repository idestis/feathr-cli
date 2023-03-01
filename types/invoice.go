package types

import "time"

// Invoice represents an invoice for a client.
type Invoice struct {
	// ID is the unique identifier of the invoice.
	ID uint `json:"id"`

	// ClientID is the unique identifier of the client associated with the invoice.
	ClientID uint `json:"client_id"`

	// Number is the invoice number.
	Number string `survey:"number" yaml:"number" json:"number"`

	// Issued is the date the invoice was issued.
	Issued time.Time `json:"issued"`

	// Due is the due date for the invoice.
	Due time.Time `survey:"due" yaml:"due" json:"due"`

	// Notes is any additional notes associated with the invoice.
	Notes string `survey:"notes" yaml:"notes" json:"notes"`

	// Items is a list of line items on the invoice.
	Items []Item `survey:"items" yaml:"items" json:"items"`
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
