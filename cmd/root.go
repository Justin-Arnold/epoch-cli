/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	//"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config struct {
	DefaultSessionDuration int `mapstructure:"default_session_duration"`
	DefaultBreakDuration   int `mapstructure:"default_break_duration"`
}

var (
	cfgFile string
	config  Config
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "epoch-cli",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.epoch-cli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.epoch-cli.yaml)")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("toml")
		viper.SetConfigName(".epoch-cli")
	}

	viper.AutomaticEnv()

	// Set default values before reading or creating the config file
	viper.SetDefault("default_session_duration", 25)
	viper.SetDefault("default_break_duration", 5)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; create it with default values
			err = viper.SafeWriteConfig()
			if err != nil {
				fmt.Println("Error creating config file:", err)
				os.Exit(1)
			}
			fmt.Println("Created new config file:", viper.ConfigFileUsed())
		} else {
			// Config file was found but another error was produced
			fmt.Println("Error reading config file:", err)
			os.Exit(1)
		}
	} else {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	// Unmarshal the config into the Config struct
	err := viper.Unmarshal(&config)
	cobra.CheckErr(err)

	// Write the config to ensure all default values are saved
	err = viper.WriteConfig()
	if err != nil {
		fmt.Println("Error writing config file:", err)
		os.Exit(1)
	}
}
