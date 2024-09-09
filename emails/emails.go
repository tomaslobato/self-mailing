package emails

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"gopkg.in/gomail.v2"
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
	if fromName == "" {
		return fmt.Errorf("FROM_NAME variable not found\n Use self-mailing setenv FROM_NAME 'Your Name' to add it")
	}
	fromAddr := os.Getenv("FROM_ADDRESS")
	if fromName == "" {
		return fmt.Errorf("FROM_ADDRESS variable not found\n Use self-mailing setenv FROM_ADDRESS 'example@email.com' to add it")
	}

	unsuscribeLink := os.Getenv("UNSUSCRIBE_LINK")
	var unsuscribeMsg string
	if unsuscribeLink == "" {
		unsuscribeMsg = ""
	} else {
		unsuscribeMsg = fmt.Sprintf(`<p>If you wish to unsuscribe, <a href=%s>Click here</a></p>`, unsuscribeLink)
	}

	if service == "sendgrid" {
		apiKey := os.Getenv("SENDGRID_KEY")
		if apiKey == "" {
			return fmt.Errorf("sendgrid key not found")
		}

		client := sendgrid.NewSendClient(apiKey)

		return sendWithSendgrid(client, fromName, fromAddr, filepath, subject, unsuscribeMsg, &emails)
	}

	if service == "gmail" {
		gmailAppPwd := os.Getenv("GMAIL_APP_PASSWORD")
		if gmailAppPwd == "" {
			return fmt.Errorf("gmail app password not found")
		}

		dialer := gomail.NewDialer("smtp.gmail.com", 587, fromAddr, gmailAppPwd)
		dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

		return sendWithGmail(dialer, fromAddr, filepath, subject, unsuscribeMsg, &emails)
	}

	return fmt.Errorf("unsupported service: %s", service)
}

func sendWithGmail(dialer *gomail.Dialer, fromAddr, filepath, subject, unsuscribeMsg string, emails *[]string) error {
	fileContent, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(*emails))
	for _, to := range *emails {
		wg.Add(1)

		go func(toAddr string) {
			defer wg.Done()

			m := gomail.NewMessage()
			m.SetHeader("From", fromAddr)
			m.SetHeader("To", toAddr)
			m.SetHeader("Subject", subject)
			m.SetBody("text/html", string(fileContent)+unsuscribeMsg)

			err = dialer.DialAndSend(m)

			if err != nil {
				errChan <- fmt.Errorf("error sending email to %s: %w", toAddr, err)
			}

			fmt.Printf("Email sent to %s\n", toAddr)
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

func sendWithSendgrid(client *sendgrid.Client, fromName, fromAddr, filepath, subject, unsuscribeMsg string, emails *[]string) error {
	fileContent, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(*emails))
	for _, to := range *emails {
		wg.Add(1)
		go func(toAddr string) {
			defer wg.Done()

			from := mail.NewEmail(fromName, fromAddr)
			to := mail.NewEmail("Recipient", toAddr)
			msg := mail.NewSingleEmail(from, subject, to, "", "")

			p := mail.NewPersonalization()
			p.AddTos(to)
			msg.AddPersonalizations(p)

			c := mail.NewContent("text/html", string(fileContent)+unsuscribeMsg)
			msg.AddContent(c)

			response, err := client.Send(msg)
			if err != nil || response.StatusCode > 300 {
				errChan <- fmt.Errorf("error sending email to %s: %s", toAddr, response.Body)
			}

			fmt.Printf("Email sent to %s - Status code: %d\n", toAddr, response.StatusCode)
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

func GetEmails(filename string) ([]string, error) {
	if filename == "" {
		return nil, fmt.Errorf("email list path not defined")
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
