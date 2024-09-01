package emails

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendEmails(listpath, filepath, subject, service string) error {
	err := godotenv.Load(".env.local")
	if err != nil {
		return fmt.Errorf("couldn't load .env.local")
	}

	emails, err := GetEmails(listpath)
	if err != nil {
		return fmt.Errorf("error getting email list: %w", err)
	}

	fromName := os.Getenv("FROM_NAME")
	fromAddr := os.Getenv("FROM_ADDRESS")

	if service == "sendgrid" {
		apiKey := os.Getenv("SENDGRID_KEY")

		fmt.Printf("From Name: %s\n", fromName)
		fmt.Printf("From Address: %s\n", fromAddr)
		fmt.Printf("SendGrid API Key: %s\n", apiKey)

		if apiKey == "" {
			return fmt.Errorf("sendgrid key not defined")
		}
		return SendWithSendgrid(apiKey, fromName, fromAddr, filepath, subject, emails)
	}

	return fmt.Errorf("unsupported service: %s", service)
}

func SendWithSendgrid(apiKey, fromName, fromAddr, filepath, subject string, emails []string) error {
	fileContent, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	client := sendgrid.NewSendClient(apiKey)

	var wg sync.WaitGroup
	errChan := make(chan error, len(emails))
	for _, to := range emails {
		wg.Add(1)
		go func(toAddr string) {
			defer wg.Done()
			err := sendEmail(client, fromName, fromAddr, toAddr, subject, fileContent)
			if err != nil {
				errChan <- fmt.Errorf("Error sending email to %s: %w\n", toAddr, err)
			}
		}(to)
	}

	wg.Wait()
	close(errChan)

	var errStrings []string
	for err := range errChan {
		errStrings = append(errStrings, err.Error())
	}

	if len(errStrings) > 0 {
		return fmt.Errorf("errors occurred while sending emails: %v", errStrings)
	}

	return nil
}

func sendEmail(client *sendgrid.Client, fromName, fromAddr, toAddr, subject string, fileContent []byte) error {
	from := mail.NewEmail(fromName, fromAddr)
	to := mail.NewEmail("Recipient", toAddr)
	msg := mail.NewSingleEmail(from, subject, to, "", "")

	p := mail.NewPersonalization()
	p.AddTos(to)
	msg.AddPersonalizations(p)

	unsuscribeMsg := `<p>if you wish to unsuscribe, <a href="pdx.vercel.app/unsuscribe">Click here</a></p>`
	c := mail.NewContent("text/html", string(fileContent)+unsuscribeMsg)
	msg.AddContent(c)

	msg.SetHeader("List-Unsubscribe", "<mailto:unsubscribe@yourdomain.com>")
	msg.SetHeader("Precedence", "Bulk")
	msg.SetHeader("X-Auto-Response-Suppress", "OOF, AutoReply")

	response, err := client.Send(msg)
	if err != nil {
		return err
	}

	fmt.Printf("Email sent to %s - Status code: %d\n", toAddr, response.StatusCode)
	fmt.Printf("Response Body: %s\n", response.Body)
	fmt.Printf("Response Headers: %v\n", response.Headers)
	return nil
}

func GetEmails(filename string) ([]string, error) {
	if filename == "" {
		return nil, fmt.Errorf("Email list path not defined")
	}

	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var emails []string
	err = json.Unmarshal(file, &emails)
	if err != nil {
		return nil, err
	}

	return emails, nil
}
