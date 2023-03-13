package cmd

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/idestis/feathr-cli/types"
	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Have a new client? Create a new client profile",
	Long: `Company profiles is a way to store the client information in a file or a SQLite database.
This information is used to generate the invoices and keep track of invoices per a client.
Simply answer the questions and the client profile will be created.`,
	Run: func(cmd *cobra.Command, args []string) {
		client := types.Client{}
		profile := types.Profile{}
		// Read the profile information
		if err := profile.Load(dataDir); err != nil {
			// FIXME: Handle the error with WARN instead of FATAL
			cobra.CheckErr(err)
		}
		info := []*survey.Question{
			{
				Name: "name",
				Prompt: &survey.Input{
					Message: "What is the client name? (e.g. Alphabet Inc.)",
					Help:    "The client name is used to identify the client in the invoice.",
				},
				Validate:  survey.Required,
				Transform: survey.Title,
			},
			{
				Name: "address",
				Prompt: &survey.Input{
					Message: "What is the client address?",
					Help:    "The physical address of the client. This is used to generate the invoice.",
				},
				Validate: survey.Required,
			},
			{
				Name: "currency",
				Prompt: &survey.Input{
					Message: "Which currency used to bill a client?",
					Help:    "The currency may differ from client to client, but often you can use the same currency for all clients.",
					Default: profile.Currency,
				},
			},
			{
				Name: "iban",
				Prompt: &survey.Input{
					Message: "What is the client IBAN?",
					Help:    "The IBAN is used to generate the invoice. If you don't know the IBAN, you can leave this field empty.",
				},
			},
			{
				Name: "bank",
				Prompt: &survey.Input{
					Message: "What is the client bank details?",
					Help:    "The bank details may be helpful tracking your income. If you don't know the bank details, you can leave this field empty.",
				},
			},
		}
		if err := survey.Ask(info, &client); err != nil {
			cobra.CheckErr(err)
		}

		var emails string
		if err := survey.AskOne(&survey.Multiline{
			Message: "What is the client emails? (starting each on separate line)",
		}, &emails); err != nil {
			cobra.CheckErr(err)
		} else {
			client.Email = []string{}
			for _, email := range strings.Split(emails, "\n") {
				client.Email = append(client.Email, strings.TrimSpace(email))
			}
		}
		fmt.Printf("%+v", client)
		if err := client.WriteClientInfo(dataDir); err != nil {
			cobra.CheckErr(err)
		}
	},
}

func init() {
	clientsCmd.AddCommand(newCmd)
}
