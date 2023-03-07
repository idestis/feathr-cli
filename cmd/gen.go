package cmd

import (
	"errors"

	"github.com/AlecAivazis/survey/v2"
	"github.com/idestis/feathr-cli/helpers"
	"github.com/idestis/feathr-cli/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	invoice int
	client  int
)

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate a PDF invoice from a prepared invoice",
	Long:  `TBD`,
	Run: func(cmd *cobra.Command, args []string) {
		storageType := viper.GetString("storage")
		client_ids, err := types.GetClientIDs(dataDir, storageType)
		cobra.CheckErr(err)
		if client > 0 {
			if !helpers.ContainsIntInSlice(client, client_ids) {
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
			client = clients[selected]
		}
	},
}

func init() {
	genCmd.Flags().IntVarP(&invoice, "invoice", "i", 0, "The invoice number to generate")
	genCmd.Flags().IntVarP(&client, "client", "c", 0, "The client name to generate the invoice for")
	invoicesCmd.AddCommand(genCmd)
}
