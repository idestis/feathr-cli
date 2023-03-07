package cmd

import (
	"errors"
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
	dataDir string
	cfgFile string
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
		// allowCmds contains a list of commands that can be run without initializing the CLI

		allowCmds := []string{"init", "version", "help"}
		_, err := os.Stat(cfgFile)
		if os.IsNotExist(err) {
			if !helpers.ContainsInSlice(cmd.Name(), allowCmds) {
				return fmt.Errorf("please run '%v init' to initialize the Feathr", appName)
			}
		} else {

		}
		return nil
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())

}

func init() {
	cobra.OnInitialize(initConfig)

}

func initConfig() {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	viper.AddConfigPath(home)
	viper.SetConfigType("yaml")
	viper.SetConfigName(".feathr-cli")
	if err := viper.ReadInConfig(); err != nil {
		cobra.CheckErr(errors.New(fmt.Sprintf("error reading config file: %v", err)))
	}
	// filepath.Join(home, fmt.Sprintf(".%v", appName))
	dataDir = filepath.Join(home, fmt.Sprintf(".%v", appName))
	cfgFile = viper.ConfigFileUsed()
}
