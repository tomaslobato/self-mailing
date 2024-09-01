package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tomaslobato/self-mailing/emails"
)

var sendCmd = &cobra.Command{
	Use:   "send [file] to [list] subject [subject]",
	Short: "Send a file as email to your list",
	Args:  cobra.ExactArgs(5),
	Run: func(cmd *cobra.Command, args []string) {
		filepath := args[0]
		listpath := args[2]
		subject := args[4]

		useSendgrid, err := cmd.Flags().GetBool("sendgrid")
		var service string
		if useSendgrid {
			service = "sendgrid"
		}

		err = emails.SendEmails(listpath, filepath, subject, service)
		if err != nil {
			fmt.Printf("Error sending emails: %v\n", err)
		}
	},
}

func init() {
	sendCmd.Flags().Bool("sendgrid", false, "Use SendGrid as the email service")
}
