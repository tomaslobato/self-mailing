package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestSendCmd(t *testing.T) {
	cmd := &cobra.Command{Use: "root"}
	cmd.AddCommand(sendCmd)

	// Test with --sendgrid flag
	args := []string{"send", "file.txt", "to", "list.json", "subject", "Test", "--sendgrid"}
	cmd.SetArgs(args)
	if err := cmd.Execute(); err != nil {
		t.Errorf("Execution failed: %v", err)
	}

	// Verify that the sendgrid flag is defined
	if sendCmd.Flags().Lookup("sendgrid") == nil {
		t.Error("sendgrid flag is not defined")
	}
}
