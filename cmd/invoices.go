package cmd

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	thousands "github.com/floscodes/golang-thousands"
	"github.com/idestis/feathr-cli/helpers"
	"github.com/idestis/feathr-cli/types"
	"github.com/spf13/cobra"
)

var hideSent bool
var clientID int

// invoicesCmd represents the invoices command
var invoicesCmd = &cobra.Command{
	Use:   "invoices",
	Args:  cobra.MatchAll(cobra.OnlyValidArgs, cobra.MaximumNArgs(1), cobra.MinimumNArgs(0)),
	Short: "Manage an invoices in a simple way",
	Long: `This command allows you to manage your invoices.

Just image days when you had to send an invoice for a client, or you had to duplicate this operation montly, weekly or even daily.
Using this command you can edit, duplicate, send and delete invoices in a simple way.`,
	Aliases: []string{"invoice", "inv"},
	Run: func(cmd *cobra.Command, args []string) {
		client_ids, err := types.GetClientIDs(dataDir)
		cobra.CheckErr(err)
		// Break early if there are no clients
		if len(client_ids) == 0 {
			cobra.CheckErr(errors.New("no clients found"))
		}
		client_names := make(map[int]types.Client, 0)
		for _, id := range client_ids {
			client, err := types.ReadClientInfo(id, dataDir)
			cobra.CheckErr(err)
			if _, ok := client_names[id]; !ok {
				client_names[id] = client
			}
			continue
		}
		invoices, _, err := types.GetInvoicesID(client_ids, dataDir)
		invoiceData := make(map[int]types.Invoice, 0)
		for id, path := range invoices {
			if _, exists := invoiceData[id]; !exists {
				invoice, err := types.ReadInvoice(id, path, dataDir)
				cobra.CheckErr(err)
				invoiceData[id] = invoice
			}
			continue
		}
		cobra.CheckErr(err)
		if len(args) > 0 {
			clientID, _ = strconv.Atoi(args[0])
		}
		// Create an options for the user to select
		options := make([]string, 0)
		for id, invoice := range invoiceData {
			if clientID > 0 && clientID != int(invoice.ClientID) {
				// Skip invoices for other clients
				continue
			}
			if hideSent && !invoice.Sent.IsZero() {
				// Hide sent invoices
				continue
			}
			total, _ := thousands.Separate(invoice.GetInvoiceTotal(), "en")
			// TODO: Add currency and thousands separator
			options = append(options, fmt.Sprintf("%d:\t%s // %v", id, client_names[int(invoice.ClientID)].Name, total))
		}
		if len(options) == 0 {
			cobra.CheckErr(errors.New("no invoices found"))
		}
		sort.Slice(options, func(i, j int) bool {
			name := strings.Split(options[i], "\t")[0]
			next := strings.Split(options[j], "\t")[0]
			return name > next
		})

		var selected string
		if err := survey.AskOne(&survey.Select{
			Message:  "Please select invoice?",
			Options:  options,
			PageSize: 10,
		}, &selected); err != nil {
			cobra.CheckErr(err)
		}

		idx, _ := strconv.Atoi(strings.Split(selected, ":")[0])
		invoice := invoiceData[idx]
		invoice.Print()
		profile := types.Profile{}
		profile.Load(dataDir)
		invoiceActions(invoice, client_names[int(invoice.ClientID)], profile)
	},
}

func init() {
	invoicesCmd.PersistentFlags().IntVarP(&clientID, "client", "c", 0, "The client ID for which you want to see or create the invoices.")
	invoicesCmd.Flags().BoolVarP(&hideSent, "no-sent", "", false, "If you want to hide invoices that have already been sent.")
	rootCmd.AddCommand(invoicesCmd)
}

func invoiceActions(invoice types.Invoice, client types.Client, profile types.Profile) {
	var action string
	if err := survey.AskOne(&survey.Select{
		Message:  "What would you like to do?",
		Options:  []string{"Edit", "Send", "Generate", "Delete", "Duplicate"},
		PageSize: 10,
	}, &action); err != nil {
		cobra.CheckErr(err)
	}
	switch action {
	case "Edit":
		fmt.Println("I am so sorry, but this feature is not yet implemented.")
	case "Send":
		fmt.Println("Send")
	case "Delete":
		var confirm bool
		if err := survey.AskOne(&survey.Confirm{
			Message: "Are you sure you want to delete this invoice?",
		}, &confirm); err != nil {
			cobra.CheckErr(err)
		}
		if !confirm {
			fmt.Println("Invoice has not been deleted.")
			return
		}
		invoice.Delete(dataDir)
		fmt.Printf("Invoice with number %d has been successfully deleted.\n", invoice.ID)
	case "Duplicate":
		clients, _ := types.GetClientIDs(dataDir)
		_, ids, _ := types.GetInvoicesID(clients, dataDir)
		last, _ := helpers.FindMaxInt(ids)
		invoice.ID = uint(last + 1)
		invoice.Number = fmt.Sprint(invoice.ID)
		if err := invoice.WriteInvoice(dataDir); err != nil {
			cobra.CheckErr(err)
		}
		fmt.Printf("Invoice with number %d has been successfully duplicated.\n", invoice.ID)
	case "Generate":
		err, path := invoice.GeneratePDF(client, profile, dataDir)
		cobra.CheckErr(err)
		fmt.Printf("Invoice generated successfully. You can find the generated file at: %s\n", path)
	}
}
