package cmd

import (
	"github.com/idestis/feathr-cli/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage feather-cli configuration after initial setup",
	ValidArgs: []string{
		"smtp",
	},
	Args: cobra.MatchAll(cobra.OnlyValidArgs, cobra.MaximumNArgs(1), cobra.MinimumNArgs(1)),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "smtp":
			cfg := types.Config{}
			err := viper.Unmarshal(&cfg)
			cobra.CheckErr(err)
			cfg.SMTP, err = types.PromptSMTP()
			cobra.CheckErr(err)
			viper.Set("smtp", cfg.SMTP)
			viper.WriteConfig()
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
