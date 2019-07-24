package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "displays version information",
	Long:  `displays version information`,
	Run:   version,
}

func version(cmd *cobra.Command, args []string) {
	fmt.Printf("cumulus version %s-%s\n", Version, Build)
}
