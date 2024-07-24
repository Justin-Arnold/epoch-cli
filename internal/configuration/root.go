package configuration

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type ConfigOptionKey string
type ConfigPartial = map[ConfigOptionKey]any

var allConfigOptions []ConfigPartial

func RegisterConfigOptions(configOptions ConfigPartial) {
	allConfigOptions = append(allConfigOptions, configOptions)
}

func init() {
	createConfig()
	setAllDefaultOptions()
	loadConfig()
	err := viper.WriteConfig()
	if err != nil {
		fmt.Println("Error writing config file:", err)
		os.Exit(1)
	}
}

func setAllDefaultOptions() {
	for _, configOptions := range allConfigOptions {
		setDefaultOptions((configOptions))
	}
}

func setDefaultOptions(options ConfigPartial) {
	for key, value := range options {
		viper.SetDefault(string(key), value)
	}
}

func createConfig() {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	viper.AddConfigPath(home)
	viper.SetConfigType("toml")
	viper.SetConfigName(".epoch-cli")
}

func loadConfig() {
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
