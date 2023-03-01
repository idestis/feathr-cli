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
	Long:  `TBD`,
	Run: func(cmd *cobra.Command, args []string) {
		config := types.Config{}
		_, err := os.Stat(cfgFile)
		if !os.IsNotExist(err) {
			fmt.Println("[!] Configuration file already exists.")
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
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
