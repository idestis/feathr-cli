package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/AlecAivazis/survey/v2"
	"github.com/idestis/feathr-cli/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// profileCmd represents the settings command
var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Configure user profile",
	Long:  `TBD`,
	ValidArgs: []string{
		"name",
		"currency",
		"due",
		"address",
		"iban",
		"bank",
	},
	Args: cobra.MatchAll(cobra.OnlyValidArgs, cobra.MaximumNArgs(1), cobra.MinimumNArgs(0)),
	Run: func(cmd *cobra.Command, args []string) {
		profile := types.Profile{}
		storageType := viper.GetString("storage")
		if !profile.Exists(dataDir, storageType) && len(args) == 0 {
			profile.Name, _ = namePrompt()
			profile.Currency, _ = currencyPrompt()
			profile.Due, _ = duePrompt()
			profile.Address, _ = addressPrompt()
			profile.IBAN, _ = ibanPrompt()
			profile.Bank, _ = bankDetailsPrompt()
			profile.Write(dataDir, storageType)
		} else if len(args) == 1 {
			profile.Load(dataDir, storageType)
			switch args[0] {
			case "name":
				profile.Name, _ = namePrompt()
				profile.Write(dataDir, storageType)
			case "currency":
				profile.Currency, _ = currencyPrompt()
				profile.Write(dataDir, storageType)
			case "due":
				profile.Due, _ = duePrompt()
				profile.Write(dataDir, storageType)
			case "address":
				profile.Address, _ = addressPrompt()
				profile.Write(dataDir, storageType)
			case "iban":
				profile.IBAN, _ = ibanPrompt()
				profile.Write(dataDir, storageType)
			case "bank":
				profile.Bank, _ = bankDetailsPrompt()
				profile.Write(dataDir, storageType)
			default:
				fmt.Println("Invalid argument")
			}
			printProfile(profile)
		} else {
			profile.Load(dataDir, storageType)
			printProfile(profile)
		}

	},
}

func init() {
	rootCmd.AddCommand(profileCmd)
}

func namePrompt() (string, error) {
	name := ""
	prompt := &survey.Input{
		Message: "What is the business name? (e.g. PE John Doe)",
	}
	err := survey.AskOne(prompt, &name, survey.WithValidator(survey.Required))
	return name, err
}

func currencyPrompt() (string, error) {
	currency := ""
	prompt := &survey.Select{
		Message: "What is the prefered currency?",
		Options: []string{"USD", "EUR", "GBP"},
		Default: "USD",
	}
	err := survey.AskOne(prompt, &currency, survey.WithValidator(survey.Required))
	return currency, err
}

func duePrompt() (string, error) {
	due := ""
	prompt := &survey.Select{
		Message: "What is the default due date?",
		Options: []string{"7", "14", "30"},
		Default: "7",
	}
	err := survey.AskOne(prompt, &due, survey.WithValidator(survey.Required))
	return due, err
}

func addressPrompt() (string, error) {
	address := ""
	prompt := &survey.Input{
		Message: "What is your address / contact information?",
	}
	err := survey.AskOne(prompt, &address, survey.WithValidator(survey.Required))
	return address, err
}

func ibanPrompt() (string, error) {
	iban := ""
	prompt := &survey.Input{
		Message: "What is your IBAN?",
	}
	err := survey.AskOne(prompt, &iban, survey.WithValidator(survey.Required))
	return iban, err
}

func bankDetailsPrompt() (string, error) {
	bankDetails := ""
	prompt := &survey.Multiline{
		Message: "What are your bank details (including correspondent bank)?",
	}
	err := survey.AskOne(prompt, &bankDetails, survey.WithValidator(survey.Required))
	return bankDetails, err
}

func profileFileExists(dataDir string) bool {
	_, err := os.Stat(filepath.Join(dataDir, "profile.json"))
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		// Handle other error cases here if necessary
	}
	return true
}

func printProfile(p types.Profile) {
	v := reflect.ValueOf(p)
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fmt.Printf("\033[32m%s\033[0m: %v\n", v.Type().Field(i).Name, field.Interface())
	}
}
