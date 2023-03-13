package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	appName = "feathr-cli"
)

var (
	dataDir string
	cfgFile string
	version Version
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "feathr-cli",
	Short: "Feather-cli provides a simple and efficient payment processing experience",
	Long: `Feathr is a Command Line Interface (CLI) tool designed to help users generate and send invoices to
multiple clients directly from their local machine's terminal. 

This tool is developed to provide an efficient and straightforward invoicing solution
for small business owners and freelancers who need to create and manage invoices for their
clients.`,

	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		viper.SetConfigName(".feathr-cli")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(home)
		err = viper.ReadInConfig()
		if err != nil {
			// Check if config file already exists
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				// Config file doesn't exist, check if command is allowed
				if cmd.Name() != "init" && cmd.Name() != "version" {
					return fmt.Errorf("please run '%v init' to initialize the Feathr CLI", appName)
				}
				return nil
			}
			// Config file exists but failed to read it, report error
			return fmt.Errorf("failed to read config file: %v", err)
		}

		return nil
	},
}

func Execute(v Version) {
	version = v
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	dataDir = filepath.Join(home, fmt.Sprintf(".%v", appName))
	cfgFile = viper.ConfigFileUsed()
}
