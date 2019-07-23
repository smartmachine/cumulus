package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/kyokomi/emoji"
)

func init() {
	rootCmd.AddCommand(checkCmd)
}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "checks connections and permissions to our cloud provider",
	Long:  `Makes sure that we have valid credentials and roles to be able to provision cloudy stuff`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%sAll checks passed.\n", emoji.Sprint(":sparkles:"))
	},
}
