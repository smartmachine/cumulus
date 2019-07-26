package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	console "go.smartmachine.io/cumulus/pkg/repl"
	"os"
	"os/user"
)

func init() {
	rootCmd.AddCommand(replCmd)
}

var replCmd = &cobra.Command{
	Use:   "repl",
	Short: "cumulus repl",
	Long:  `Type your commands at the prompt and test parsing.`,
	Run:   repl,
}

func repl(cmd *cobra.Command, args []string) {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is the Cumulus programming language!\n", usr.Username)
	fmt.Println("Feel free to type in commands")
	console.Start(os.Stdin, os.Stdout)
}
