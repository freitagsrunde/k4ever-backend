package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Create the cobra root command for proper cli usage
var rootCmd = &cobra.Command{
	Use:   "k4ever",
	Short: "k4ever is a shopping site",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

// Handles all global flags
func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringP("version", "", "1.0.0", "The programs current version")
	viper.BindPFlag("version", rootCmd.PersistentFlags().Lookup("version"))
	viper.SetDefault("version", "0.0.1")
}

func initConfig() {
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
