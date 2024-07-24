package configuration

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type ConfigOption string

const (
	// SessionDuration is the key for the default session duration option
	DefaultSessionDuration string = "default_session_duration"
	DefaultBreakDuration   string = "default_break_duration"
)

// DefaultValues stores the default values for each option
var DefaultConfigOptionValues = map[string]interface{}{
	DefaultSessionDuration: 25,
	DefaultBreakDuration:   5,
}

func CreateConfig() {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	viper.AddConfigPath(home)
	viper.SetConfigType("toml")
	viper.SetConfigName(".epoch-cli")
}

func SetDefaultOptions() {
	for key, value := range DefaultConfigOptionValues {
		viper.SetDefault(key, value)
	}
}

func LoadConfig() {
	err := viper.ReadInConfig()
	if err == nil {
		printConfigUsed()
		return
	}
	if isConfigNotFoundError(err) {
		createNewConfig()
		return
	}
	handleConfigReadError(err)
}

func printConfigUsed() {
	fmt.Println("Using config file:", viper.ConfigFileUsed())
}

func isConfigNotFoundError(err error) bool {
	_, ok := err.(viper.ConfigFileNotFoundError)
	return ok
}

func createNewConfig() {
	err := viper.SafeWriteConfig()
	if err != nil {
		fmt.Println("Error creating config file:", err)
		os.Exit(1)
	}
	fmt.Println("Created new config file:", viper.ConfigFileUsed())
}

func handleConfigReadError(err error) {
	fmt.Println("Error reading config file:", err)
	os.Exit(1)
}
