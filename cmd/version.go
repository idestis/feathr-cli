package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

type Version struct {
	Date    string
	Version string
	Commit  string
}

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Feathr",
	Long:  `All software has versions. This is Feathr's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s (%s) %s\n", version.Version, version.Commit, version.Date)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
