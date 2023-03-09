package cmd

import (
	"strings"

	"github.com/cheynewallace/tabby"
	"github.com/idestis/feathr-cli/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var search string

// clientsCmd represents the clients command
var clientsCmd = &cobra.Command{
	Use:     "clients",
	Short:   "Manage the clients in a simple way",
	Args:    cobra.NoArgs,
	Aliases: []string{"client", "clt"},
	Long:    `TBD`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get the storage type from the config file.
		storageType := viper.GetString("storage")

		// Get the IDs of all clients.
		ids, err := types.GetClientIDs(dataDir, storageType)
		cobra.CheckErr(err)

		// Create a new tabby table to display the clients.
		t := tabby.New()
		t.AddHeader("ID", "Name", "Address", "Emails", "IBAN", "Bank")

		// Loop through each client and add them to the table if they match the search query.
		for _, id := range ids {
			// Read the client information from the file or SQLite database.
			client, err := types.ReadClientInfo(id, dataDir, storageType)
			cobra.CheckErr(err)

			// Check each field of the client to see if it matches the search query.
			matches := false
			for _, field := range []string{client.Name, client.IBAN} {
				if !strings.Contains(strings.ToLower(field), strings.ToLower(search)) {
					continue
				}
				matches = true
			}

			// Add the client to the table if it matches the search query.
			if matches {
				t.AddLine(client.ID, client.Name, client.Address, strings.Join(client.Email, ", "), client.IBAN, client.Bank)
			}
		}

		// Print the table.
		t.Print()
	},
}

func init() {
	rootCmd.AddCommand(clientsCmd)
	clientsCmd.Flags().StringVarP(&search, "search", "s", "", "Search for a client using name or IBAN")
}
