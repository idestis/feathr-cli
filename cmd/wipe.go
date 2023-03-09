package cmd

import (
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

var (
	profile bool
	data    bool
)

// wipeCmd represents the wipe command
var wipeCmd = &cobra.Command{
	Use:   "wipe",
	Short: "A simple wipe command to destroy stored data",
	Args:  cobra.NoArgs,
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

	},
}

func init() {
	rootCmd.AddCommand(wipeCmd)
	wipeCmd.Flags().BoolVarP(&profile, "skip-profile", "p", false, "Do not wipe user profile settings")
	wipeCmd.Flags().BoolVarP(&data, "skip-data", "d", false, "Do not wipe data")
}
