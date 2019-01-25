package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the version number",
	Long: `Well, it just prints the version number
			nothing else really`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Version 0.0.1 this needs to be made dynamic")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
