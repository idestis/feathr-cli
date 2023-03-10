package cmd

import (
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

var (
	profile     bool
	data        bool
	application bool
)

// wipeCmd represents the wipe command
var wipeCmd = &cobra.Command{
	Use:   "wipe",
	Short: "A simple wipe command to destroy stored data",
	Long: `This command will wipe all data stored by the application.
It will ask for confirmation before wiping the data and have granualar flags to skip wiping certain data.

By default, it will wipe all data except the application configuration, which is stored in the config file.
Please note, that configuration stores some sensitive information, such as the SMTP server credentials.

You can use the --skip-profile and --skip-data flags to skip wiping the user profile settings and user data respectively.
	`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		var confirm bool
		err := survey.AskOne(&survey.Confirm{
			Message: "Are you sure you want to wipe all data?",
		}, &confirm)
		if err != nil {
			cobra.CheckErr(err)
		}
		if !confirm {
			return
		}
		if !profile {
			// Wipe user profile settings
			os.RemoveAll(filepath.Join(dataDir, "profile.json"))
		}

		if !data {
			// Wipe data
			os.RemoveAll(filepath.Join(dataDir, "data"))
		}

		if !application {
			// Wipe application configuration
			os.RemoveAll(cfgFile)
		}

	},
}

func init() {
	rootCmd.AddCommand(wipeCmd)
	wipeCmd.Flags().BoolVarP(&profile, "skip-profile", "", false, "Do not wipe user profile settings")
	wipeCmd.Flags().BoolVarP(&data, "skip-data", "", false, "Do not wipe user data")
	wipeCmd.Flags().BoolVarP(&application, "application-data", "", true, "Do not wipe application configuration")
}
