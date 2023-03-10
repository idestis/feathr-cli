package cmd

import (
	"fmt"
	"os"

	"github.com/idestis/feathr-cli/types"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialization wizard for Feathr",
	Long: `The init command is a configuration wizard that guides users through the initial setup process for the feathr-cli command-line interface. 

When a user runs this command, they will be presented with a series of prompts that allow them to customize
their experience with a CLI. The prompt will ask the user to select a storage type for the CLI. 
They will be able to choose between using a file system or a SQLite database to store their data. 
Depending on the storage type selected, the wizard will generate a default catalog structure appropriate for the storage type.
Additionally, the wizard will ask the user if they would like to generate a PDF version of the invoice automatically upon creation.`,
	Run: func(cmd *cobra.Command, args []string) {
		config := types.Config{}
		_, err := os.Stat(cfgFile)
		if !os.IsNotExist(err) {
			fmt.Println("Configuration file already exists.")
			os.Exit(1)
		}
		if err := survey.AskOne(&survey.Confirm{
			Message: "Would you like to generate a PDF version of the invoice automatically upon creation?",
			Default: true,
		}, &config.GenOnCreate); err != nil {
			cobra.CheckErr(err)
		}

		if err := survey.AskOne(&survey.Confirm{
			Message: "Would you like to generate a PDF version of the invoice automatically upon update?",
			Default: true,
		}, &config.GenOnUpdate); err != nil {
			cobra.CheckErr(err)
		}

		var setupSMPT bool
		if err := survey.AskOne(&survey.Confirm{
			Message: "Would you like to setup SMTP settings to send emails?",
			Default: false,
		}, &setupSMPT); err != nil {
			cobra.CheckErr(err)
		}
		viper.Set("gen_on_create", config.GenOnCreate)
		viper.Set("gen_on_update", config.GenOnUpdate)
		if !setupSMPT {
			viper.Set("smtp", nil)
			fmt.Println("Skipping SMTP setup. You can configure SMTP settings later by running 'feathr-cli config smtp'.")
		} else {
			config.SMTP, err = types.PromptSMTP()
			cobra.CheckErr(err)
			viper.Set("smtp", config.SMTP)
		}

		err = os.MkdirAll(fmt.Sprintf("%v/data/clients", dataDir), 0700)
		cobra.CheckErr(err)

		viper.SafeWriteConfig()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
