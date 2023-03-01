package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// invoicesCmd represents the invoices command
var invoicesCmd = &cobra.Command{
	Use:   "invoices",
	Short: "Manage an invoices in a simple way",
	Long:  `TBD`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("invoices called")
	},
}

func init() {
	rootCmd.AddCommand(invoicesCmd)
}
