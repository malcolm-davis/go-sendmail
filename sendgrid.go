// Wrapper to SendGrid's Go Library
// https://github.com/sendgrid/sendgrid-go
package sendmail

import (
	"fmt"
	"log"

	"github.com/malcolm-davis/go-stopwatch"
	sendgrid "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// SendMailManager provides functionality to send emails via SendGrid
type TrilloSendMail struct {
	APIKey string
	client *sendgrid.Client

	// User defined logger function.
	Logger func(string, ...interface{})
}

// NewSendGrid creates a new instance of TrilloSendMail
func NewSendGrid(apiKey string) (*TrilloSendMail, error) {
	client := sendgrid.NewSendClient(apiKey)
	if client == nil {
		return nil, fmt.Errorf("Failed to create SendGrid client")
	}

	manager := &TrilloSendMail{
		APIKey: apiKey,
		client: client,
	}

	return manager, nil
}

func (t *TrilloSendMail) SendMail(fromName, fromEmail, toName, toEmail, subject, plainTextContent, htmlContent string) (response *Response, err error) {
	timer := stopwatch.Start("SendMail", stopwatch.LogStop)
	defer func() {
		timer.StopE(err)
	}()

	from := mail.NewEmail(fromName, fromEmail)
	to := mail.NewEmail(toName, toEmail)
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	return t.post(message)
}

func (t *TrilloSendMail) SendMessage(message *Message) (response *Response, err error) {
	timer := stopwatch.Start("SendMessage", stopwatch.LogStop)
	defer func() {
		timer.StopE(err)
	}()

	err = message.Validate()
	if err != nil {
		return nil, err
	}

	from := mail.NewEmail(message.FromEmail.Name, message.FromEmail.Address)
	to := mail.NewEmail(message.Recipients[0].Name, message.Recipients[0].Address)

	// sendgrid use a single email
	email := mail.NewSingleEmail(from, message.Subject, to, message.PlainTextContent, message.HtmlContent)

	// add any attachments
	for _, attachment := range message.Attachments {
		// Create a new attach
		attach := mail.NewAttachment()
		attach.SetContent(attachment.Base64Content)
		attach.SetType(attachment.ContentType)
		attach.SetFilename(attachment.Filename)

		// "attachment" or "inline" for displaying in email body
		attach.SetDisposition("attachment")

		// Add the attachment to the message
		email.AddAttachment(attach)
	}

	return t.post(email)
}

func (t *TrilloSendMail) post(email *mail.SGMailV3) (response *Response, err error) {
	client := t.client
	if client == nil {
		client = sendgrid.NewSendClient(t.APIKey)
		if client == nil {
			return nil, fmt.Errorf("Failed to create SendGrid client")
		}
		t.logf("Created new SendGrid client")
	}

	trilloResponse, err := client.Send(email)
	if err != nil {
		return nil, err
	}

	mapString := mapStringSliceToString(trilloResponse.Headers)
	t.logf("Send email: status_code=%d, body=%s, headers=%s", trilloResponse.StatusCode, trilloResponse.Body, mapString)

	response = &Response{
		StatusCode: trilloResponse.StatusCode,
		Body:       trilloResponse.Body,
		Headers:    trilloResponse.Headers,
	}

	// Accept any 2xx response as success (SendGrid returns 202 Accepted).
	if trilloResponse.StatusCode < 200 || trilloResponse.StatusCode >= 300 {
		return response, fmt.Errorf("Failed to send email: %s", response.Body)
	}
	return response, nil
}

// logf logs message either via defined user logger or via system one if no user logger is defined.
func (t *TrilloSendMail) logf(f string, args ...interface{}) {
	if t.Logger != nil {
		t.Logger(f, args...)
	} else {
		log.Printf(f, args...)
	}
}
