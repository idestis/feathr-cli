package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	settings bool
	data     bool
)

// wipeCmd represents the wipe command
var wipeCmd = &cobra.Command{
	Use:   "wipe",
	Short: "A simple wipe command to destroy stored data",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("wipe called")
	},
}

func init() {
	rootCmd.AddCommand(wipeCmd)
	wipeCmd.Flags().BoolVarP(&settings, "no-settings", "s", true, "Do not wipe settings")
	wipeCmd.Flags().BoolVarP(&data, "no-data", "d", true, "Do not wipe data")
}
