package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// statsCmd represents the stats command
var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "A brief overview on user statistic",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("stats called")
	},
}

func init() {
	rootCmd.AddCommand(statsCmd)
}
