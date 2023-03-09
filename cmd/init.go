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
		// General configuration settings
		if err := survey.AskOne(&survey.Select{
			Message: "What storage do you want to use?",
			Options: []string{"sqlite", "file"},
			Default: "sqlite",
		}, &config.Storage); err != nil {
			cobra.CheckErr(err)
		}
		viper.Set("storage", config.Storage)
		if config.Storage == "sqlite" {
			// TODO: Create a database file
		} else {
			os.MkdirAll(fmt.Sprintf("%s/data", dataDir), 0700)
		}
		if err := survey.AskOne(&survey.Confirm{
			Message: "Would you like to generate a PDF version of the invoice automatically upon creation?",
			Default: true,
		}, &config.GenOnCreate); err != nil {
			cobra.CheckErr(err)
		} else {
			viper.Set("gen_on_create", config.GenOnCreate)
		}

		if err := survey.AskOne(&survey.Confirm{
			Message: "Would you like to generate a PDF version of the invoice automatically upon update?",
			Default: true,
		}, &config.GenOnUpdate); err != nil {
			cobra.CheckErr(err)
		} else {
			viper.Set("gen_on_update", config.GenOnUpdate)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
