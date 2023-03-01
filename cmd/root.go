package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/idestis/feathr-cli/helpers"
)

const (
	appName = "feathr-cli"
)

var (
	cfgFile string
	dataDir = filepath.Join(os.Getenv("HOME"), "."+appName)
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
		if cfgFile == "" {
			home, err := os.UserHomeDir()
			cobra.CheckErr(err)
			cfgFile = filepath.Join(home, ".feathr-cli.yaml")
		}
		// allowCmds contains a list of commands that can be run without initializing the CLI
		allowCmds := []string{"init", "version", "help"}
		viper.SetConfigType("yaml")
		viper.SetConfigFile(cfgFile)
		_, err := os.Stat(cfgFile)
		if os.IsNotExist(err) {
			if !helpers.ContainsInSlice(cmd.Name(), allowCmds) {
				return fmt.Errorf("please run '%v init' to initialize the Feathr", appName)
			}
		} else {
			if err := viper.ReadInConfig(); err != nil {
				return fmt.Errorf("error reading config file: %v", err)
			}
		}
		return nil
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())

}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/feathr-cli/config.yaml)")
}
