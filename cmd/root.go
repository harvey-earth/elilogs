package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Root is the root elilogs command
func Root() *cobra.Command {
	return rootCmd
}

var rootCmd = &cobra.Command{
	Use:   "elilogs",
	Short: "interact with elasticsearch",
	Long:  `elilogs is a CLI application that allows easy interaction with elasticsearch. This app can easily list cluster and index information and can even run multi-index queries easily.`,

	Version: "0.0.1",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "debug output")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	rootCmd.MarkFlagsMutuallyExclusive("debug", "verbose")

}

// Configure reads in configuration file and environment variables
func Configure() {
	viper.SetDefault("logLevel", "warn")

	viper.SetEnvPrefix("ELILOGS")
	viper.AutomaticEnv()
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/elilogs")
	viper.AddConfigPath("$HOME/.elilogs")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("Configuration file not found")
			os.Exit(1)
		} else {
			fmt.Println("Fatal error in configuration: ", err)
			os.Exit(1)
		}
	}
	SetLogLevel()
}

// SetLogLevel sets viper string logLevel
func SetLogLevel() {
	lvl := viper.GetString("core.log_level")
	if lvl != "" {
		viper.Set("logLevel", lvl)
	}
	if c, _ := rootCmd.Flags().GetBool("verbose"); c {
		viper.Set("logLevel", "info")
	}
	if c, _ := rootCmd.Flags().GetBool("debug"); c {
		viper.Set("logLevel", "debug")
	}
}
