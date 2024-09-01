package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "self-mailing",
	Short: "A self-hosted CLI tool for sending emails to an email list",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Please specify a command. Use self-mailing --help")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(sendCmd)
	rootCmd.AddCommand(setEnvCmd)
}
