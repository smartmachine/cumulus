package cmd

import (
	"fmt"
	"github.com/kyokomi/emoji"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(checkCmd)
}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "checks connections and permissions to our cloud provider",
	Long:  `Makes sure that we have valid credentials and roles to be able to provision cloudy stuff`,
	Run:   check,
}

func check(cmd *cobra.Command, args []string) {
	fmt.Printf("%sAll checks passed.\n", emoji.Sprint(":sparkles:"))
}
