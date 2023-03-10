package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/idestis/feathr-cli/types"
	"github.com/spf13/cobra"
)

var (
	search   string
	pageSize int
)

// clientsCmd represents the clients command
var clientsCmd = &cobra.Command{
	Use:     "clients",
	Short:   "Manage the clients in a simple way",
	Args:    cobra.NoArgs,
	Aliases: []string{"client", "clt"},
	Run: func(cmd *cobra.Command, args []string) {
		// Get the IDs of all clients.
		ids, err := types.GetClientIDs(dataDir)
		cobra.CheckErr(err)
		if len(ids) == 0 {
			cobra.CheckErr("You have not added any clients yet.\nPlease proceed with adding a client first using the 'feathr-cli clients new' command.")
		}

		// Create a new tabby table to display the clients.
		clients := make(map[string]types.Client)
		clientNames := []string{}
		// Loop through each client and add them to the table if they match the search query.
		for _, id := range ids {
			// Read the client information from the file or SQLite database.
			client, err := types.ReadClientInfo(id, dataDir)
			cobra.CheckErr(err)

			// Check each field of the client to see if it matches the search query.
			matches := false
			for _, field := range []string{client.Name} {
				if !strings.Contains(strings.ToLower(field), strings.ToLower(search)) {
					continue
				}
				matches = true
			}

			// Add the client to the table if it matches the search query.
			if matches {
				clients[client.Name] = client
				clientNames = append(clientNames, fmt.Sprintf("%v: %v", client.ID, client.Name))
			}
		}
		if len(clientNames) == 0 {
			cobra.CheckErr("No clients match the search query.")
		}
		var clientName string
		err = survey.AskOne(&survey.Select{
			Message:  "Select a client",
			Options:  clientNames,
			PageSize: pageSize,
		}, &clientName)
		cobra.CheckErr(err)
		client := clients[strings.TrimSpace(strings.Split(clientName, ":")[1])]
		client.Print()
		var action string
		err = survey.AskOne(&survey.Select{
			Message: "What do you want to do?",
			Help:    "Select an action to perform on the client",
			Options: []string{
				"Edit client",
				"Add invoice",
				"View invoices",
				"Quit",
			},
		}, &action)
		cobra.CheckErr(err)
		switch action {

		case "Edit client":
			fmt.Println("Editing the client is not avaialble yet.")
		case "Add invoice":
			createCmd.Run(cmd, []string{fmt.Sprint(client.ID)})
		case "View invoices":
			invoicesCmd.Run(cmd, []string{fmt.Sprint(client.ID)})
		case "Quit":
			fmt.Println("Bye!")
			os.Exit(0)
		}
	},
}

func init() {
	rootCmd.AddCommand(clientsCmd)
	clientsCmd.Flags().StringVarP(&search, "search", "s", "", "Search for a client using name")
	clientsCmd.Flags().IntVarP(&pageSize, "limit", "l", 10, "The number of clients to display at one page")
}
