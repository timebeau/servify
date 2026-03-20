package cli

import (
	"fmt"

	"servify/apps/server/internal/version"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of servify",
	Long:  `All software has versions. This is servify's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version: %s\nCommit: %s\nBuildTime: %s\n", version.Version, version.Commit, version.BuildTime)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
