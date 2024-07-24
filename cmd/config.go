package cmd

import (
	"fmt"
	"os"

	"github.com/Justin-Arnold/epoch-cli/internal/configuration"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configCommand = &cobra.Command{
	Use:   "config",
	Short: "Used to make changes to the config from the command line",
	Run:   changeConfigSetting,
	// Commands should be as follows
	// - session [time]
}

func changeConfigSetting(command *cobra.Command, commandLineArguments []string) {
	if len(commandLineArguments) != 2 {
		fmt.Printf("expected 2 arguments, got %d", len(commandLineArguments))

	}

	setting := commandLineArguments[0]
	value := commandLineArguments[1]

	var key string
	var err error

	switch setting {
	case "session":
		key = string(configuration.DefaultFocusDuration)
	case "break":
		key = string(configuration.DefaultBreakDuration)
	default:
		fmt.Printf("unknown setting: %s", setting)
		os.Exit(1)
	}

	// Set the new value
	viper.Set(key, value)

	// Write the changes to the config file
	err = viper.WriteConfig()
	if err != nil {
		fmt.Printf("error writing config: %w", err)
		os.Exit(1)
	}

	fmt.Printf("Updated %s to %v\n", setting, value)
}

func init() {
	rootCmd.AddCommand(configCommand)
}
