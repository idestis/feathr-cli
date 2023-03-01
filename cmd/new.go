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
	Short: "A command to add a new clients",
	Long:  `TBD`,
	Run: func(cmd *cobra.Command, args []string) {
		client := types.Client{}
		info := []*survey.Question{
			{
				Name: "name",
				Prompt: &survey.Input{
					Message: "What is the client name? (e.g. Alphabet Inc.)",
				},
				Validate:  survey.Required,
				Transform: survey.Title,
			},
			{
				Name: "address",
				Prompt: &survey.Input{
					Message: "What is the client address?",
				},
				Validate: survey.Required,
			},
			{
				Name: "iban",
				Prompt: &survey.Input{
					Message: "What is the client IBAN?",
				},
			},
			{
				Name: "bank",
				Prompt: &survey.Input{
					Message: "What is the client bank details?",
				},
			},
		}
		if err := survey.Ask(info, &client); err != nil {
			cobra.CheckErr(err)
		}

		var emails string
		if err := survey.AskOne(&survey.Multiline{
			Message: "What is the client emails? (starting each on a new line)",
		}, &emails); err != nil {
			cobra.CheckErr(err)
		} else {
			client.Email = []string{}
			for _, email := range strings.Split(emails, "\n") {
				client.Email = append(client.Email, strings.TrimSpace(email))
			}
		}
		fmt.Printf("%+v", client)
		// TODO: Placeholder on how to write the client info to a file
		if err := client.WriteClientInfo(dataDir, "file"); err != nil {
			cobra.CheckErr(err)
		}
	},
}

func init() {
	clientsCmd.AddCommand(newCmd)
}
