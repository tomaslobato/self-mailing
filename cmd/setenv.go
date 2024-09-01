package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var setEnvCmd = &cobra.Command{
	Use:   "setenv <NAME> <VALUE>",
	Short: "Set environment variables at .env.local",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		name := strings.ToUpper(args[0])
		value := args[1]

		err := appendEnvVariable(".env.local", name, value)
		if err != nil {
			fmt.Printf("Error reading .env.local", err)
			return
		}
	},
}

func appendEnvVariable(filename, name, value string) error {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		os.Create(filename)
	} else if err != nil {
		return fmt.Errorf("failed to read env file: %w", err)
	}

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open env file: %w", err)
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, "%s=%s\n", name, value)
	if err != nil {
		return fmt.Errorf("failed to write env file: %w", err)
	}

	return nil
}
