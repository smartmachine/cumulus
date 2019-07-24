package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use: "cumulus",
	Short: "Cumulus is an AWS self-deploying framework",
	Long: `Cumulus is examining AWS self-provisioning.  Eventually we want
to distribute entire platforms with a single binary distributable.`,
}

var Version string
var Build   string

func init() {
	rootCmd.PersistentFlags().String("profile", "", "Profile to use.")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
