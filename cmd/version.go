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
		fmt.Print(`Version: ` + config.Version() + "\nBranch: " + config.GitBranch() + "\nCommitHash: " + config.GitCommit() + "\nBuild time: " + config.BuildTime() + "\n")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
