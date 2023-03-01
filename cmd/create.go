package cmd

import (
	"errors"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/cheynewallace/tabby"
	"github.com/idestis/feathr-cli/helpers"
	"github.com/idestis/feathr-cli/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var clientID int

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Simply create an invoice",
	Args:  cobra.NoArgs,
	Long:  `TBD`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get the storage type
		storageType := viper.GetString("storage")
		// Read the data in order to get the client IDs
		client_ids, err := types.GetClientIDs(dataDir, storageType)
		cobra.CheckErr(err)
		if clientID > 0 {
			// Check if the client ID is valid
			if !helpers.ContainsIntInSlice(clientID, client_ids) {
				cobra.CheckErr(errors.New("the client ID is not valid"))
			}
		} else {
			// Read the client information from the file or SQLite database.
			clients := make(map[string]int)
			for _, id := range client_ids {
				client, err := types.ReadClientInfo(id, dataDir, storageType)
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
				Options: options,
			}, &selected); err != nil {
				cobra.CheckErr(err)
			}
		}

		invoice := types.Invoice{}
		another := true
		if err := survey.AskOne(&survey.Input{
			// TODO: Need to check if the invoice number is unique
			// TODO: Provide default value based on the global next invoice number
			Message: "What is the invoice number?",
		}, &invoice.Number); err != nil {
			cobra.CheckErr(err)
		}
		t := tabby.New()
		t.AddHeader("Description", "Qty.", "Price", "Total")
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
					},
					Validate: survey.Required,
				},
				{
					Name: "quantity",
					Prompt: &survey.Input{
						Message: "How many hours did you work? (Qty.)",
					},
					Validate: survey.Required,
				},
				{
					Name: "unit_price",
					Prompt: &survey.Input{
						Message: "What is the hourly rate? (Unit Price)",
					},
					Validate: survey.Required,
				},
			}
			if err := survey.Ask(item, &temp); err != nil {
				cobra.CheckErr(err)
			}
			invoice.Items = append(invoice.Items, temp)
			t.AddLine(temp.Description, temp.Quantity, temp.UnitPrice, temp.UnitPrice*temp.Quantity)
			if err := survey.AskOne(&survey.Confirm{
				Message: "Add another item?",
			}, &another); err != nil {
				cobra.CheckErr(err)
			}
		}

		if err := survey.AskOne(&survey.Multiline{
			Message: "Any additional notes to invoice?",
		}, &invoice.Notes); err != nil {
			cobra.CheckErr(err)
		}
		fmt.Printf("Invoice #%v \n\n", invoice.Number)
		t.Print()
		fmt.Printf("\nNotes: %v\n", invoice.Notes)
		// TODO: Add a confirmation step and write step
	},
}

func init() {
	invoicesCmd.AddCommand(createCmd)
	createCmd.Flags().IntVarP(
		&clientID,
		"client-id",
		"c",
		0,
		"Create an invoice for a specific client ID, otherwise you will be prompted to select a client.",
	)
}
