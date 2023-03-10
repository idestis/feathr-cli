package cmd

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/idestis/feathr-cli/helpers"
	"github.com/idestis/feathr-cli/types"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Simply create an invoice",
	Long: `Create an invoice for the client by selecting the client from the list of available clients.

When you need to create invoice for the new client, proceed with client creation first.
If you know the client ID, you can pass it as an argument or the flag --client to the command.`,
	Args: cobra.MatchAll(cobra.OnlyValidArgs, cobra.MaximumNArgs(1), cobra.MinimumNArgs(0)),
	Run: func(cmd *cobra.Command, args []string) {
		// Read the data in order to get the client IDs
		client_ids, err := types.GetClientIDs(dataDir)
		invoice := types.Invoice{}
		cobra.CheckErr(err)
		if len(args) > 0 {
			clientID, _ = strconv.Atoi(args[0])
		}
		if clientID > 0 {
			// Check if the client ID is valid
			if !helpers.ContainsIntInSlice(clientID, client_ids) {
				cobra.CheckErr(errors.New("the client ID is not valid"))
			}
			invoice.ClientID = uint(clientID)
		} else {
			// Read the client information from the file or SQLite database.
			clients := make(map[string]int)
			for _, id := range client_ids {
				client, err := types.ReadClientInfo(id, dataDir)
				cobra.CheckErr(err)
				clients[client.Name] = id
			}
			options := make([]string, 0)
			for k := range clients {
				options = append(options, k)
			}
			var selected string
			if err := survey.AskOne(&survey.Select{
				Message: "Which client is this invoice for?",
				Help: `Select the client from the list of available clients,
if you need to create invoice for the new client, 
proceed with client creation first.`,
				Options: options,
			}, &selected); err != nil {
				cobra.CheckErr(err)
			}
			invoice.ClientID = uint(clients[selected])
		}
		_, invoiceIDs, _ := types.GetInvoicesID(client_ids, dataDir)
		another := true
		if len(invoiceIDs) > 0 {
			next, _ := helpers.FindMaxInt(invoiceIDs)
			invoice.ID = uint(next + 1)
		} else {
			invoice.ID = 1
		}
		if err := survey.AskOne(&survey.Input{
			// TODO: Need to check if the invoice number is unique
			// TODO: Provide default value based on the global next invoice number
			Message: "What is the invoice number?",
			Default: fmt.Sprintf("%v", invoice.ID),
		}, &invoice.Number); err != nil {
			cobra.CheckErr(err)
		}
		now := time.Now()

		for {
			if !another {
				break
			}
			temp := types.Item{}
			item := []*survey.Question{
				{
					Name: "description",
					Prompt: &survey.Input{
						Message: "What service did you provide?",
						Suggest: helpers.SuggestService,
						Help:    "Provide a short description of the service you provided.",
					},
					Validate: survey.Required,
				},
				{
					Name: "quantity",
					Prompt: &survey.Input{
						Message: "How many hours did you work? (Qty.)",
						Help:    "Provide the number of hours you worked on the project",
						Default: "1",
					},
					// TODO: Validate that the value is a number, more than 0 and doesn't have special characters
					Validate: survey.Required,
				},
				{
					Name: "unit_price",
					Prompt: &survey.Input{
						Message: "What is the hourly rate? (Unit Price)",
						Help:    "Provide the hourly rate for the service you provided.",
					},
					// Validate if number, more than 0 and doesn't have special characters, except for the dot
					Validate: survey.Required,
				},
			}
			if err := survey.Ask(item, &temp); err != nil {
				cobra.CheckErr(err)
			}
			invoice.Items = append(invoice.Items, temp)
			if err := survey.AskOne(&survey.Confirm{
				Message: "Add another item?",
				Help:    "In case if you need to add more services to the same invoice.",
			}, &another); err != nil {
				cobra.CheckErr(err)
			}
		}

		if err := survey.AskOne(&survey.Multiline{
			Message: "Any additional notes to invoice?",
			Help:    "Provide any additional notes to the invoice.",
		}, &invoice.Notes); err != nil {
			cobra.CheckErr(err)
		}
		invoice.Issued = now
		profile := types.Profile{}
		profile.Load(dataDir)
		due, _ := strconv.Atoi(profile.Due)
		invoice.Due = now.AddDate(0, 0, due)
		invoice.Print()
		invoice.WriteInvoice(dataDir)
	},
}

func init() {
	invoicesCmd.AddCommand(createCmd)
}
