package cmd

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/idestis/feathr-cli/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var clientID int
var hideSent bool

// invoicesCmd represents the invoices command
var invoicesCmd = &cobra.Command{
	Use:     "invoices",
	Short:   "Manage an invoices in a simple way",
	Long:    `TBD`,
	Aliases: []string{"invoice"},
	Run: func(cmd *cobra.Command, args []string) {
		storageType := viper.GetString("storage")
		client_ids, err := types.GetClientIDs(dataDir, storageType)
		cobra.CheckErr(err)
		// Break early if there are no clients
		if len(client_ids) == 0 {
			cobra.CheckErr(errors.New("no clients found"))
		}
		client_names := make(map[int]types.Client, 0)
		for _, id := range client_ids {
			client, err := types.ReadClientInfo(id, dataDir, storageType)
			cobra.CheckErr(err)
			if _, ok := client_names[id]; !ok {
				client_names[id] = client
			}
			continue
		}
		invoices, _, err := types.GetInvoicesID(client_ids, dataDir, storageType)
		invoiceData := make(map[int]types.Invoice, 0)
		for id, path := range invoices {
			if _, exists := invoiceData[id]; !exists {
				invoice, err := types.ReadInvoice(id, path, dataDir, storageType)
				cobra.CheckErr(err)
				invoiceData[id] = invoice
			}
			continue
		}
		cobra.CheckErr(err)

		// Create an options for the user to select
		options := make([]string, 0)
		for id, invoice := range invoiceData {
			if clientID != 0 && clientID != id {
				// Skip invoices for other clients
				continue
			}
			if hideSent && !invoice.Sent.IsZero() {
				// Hide sent invoices
				continue
			}
			total := invoice.GetInvoiceTotal()
			// TODO: Add currency and thousands separator
			options = append(options, fmt.Sprintf("%d:\t%s\t\t%g", id, client_names[id].Name, total))
		}
		if len(options) == 0 {
			cobra.CheckErr(errors.New("no invoices found"))
		}
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
	},
}

func init() {
	invoicesCmd.PersistentFlags().IntVarP(&clientID,
		"client-id",
		"c",
		0,
		"Select the client ID")
	invoicesCmd.Flags().BoolVarP(&hideSent, "no-sent", "", false, "If you want to hide invoices that have already been sent.")
	rootCmd.AddCommand(invoicesCmd)
}
