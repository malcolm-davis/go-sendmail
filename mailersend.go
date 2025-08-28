// Wrapper to MailerSend Go Library
// https://github.com/MailerSend/MailerSend-apiv3-go
package sendmail

import (
	"context"
	"fmt"
	"log"

	"github.com/mailersend/mailersend-go"
	"github.com/malcolm-davis/go-stopwatch"
)

type MailerSend struct {
	client *mailersend.Mailersend
	token  string

	// override logging using user defined logger function.
	Logger func(string, ...interface{})
}

func NewMailerSend(apiToken string) (*MailerSend, error) {
	client := mailersend.NewMailersend(apiToken)
	if client == nil {
		return nil, fmt.Errorf("Failed to create MailerSend client")
	}

	manager := &MailerSend{
		token:  apiToken,
		client: client,
	}

	return manager, nil
}

func (ms *MailerSend) SendMail(fromName, fromEmail, toName, toEmail, subject, plainTextContent, htmlContent string) (response *Response, err error) {
	timer := stopwatch.Start("SendMail", stopwatch.LogStop)
	defer func() {
		timer.StopE(err)
	}()

	message := ms.client.Email.NewMessage()

	from := mailersend.From{
		Name:  fromName,
		Email: fromEmail,
	}

	recipients := []mailersend.Recipient{
		{
			Name:  toName,
			Email: toEmail,
		},
	}

	message.SetFrom(from)
	message.SetRecipients(recipients)
	message.SetSubject(subject)
	message.SetHTML(htmlContent)
	message.SetText(plainTextContent)

	return ms.post(message)
}

func (ms *MailerSend) SendMessage(message *Message) (response *Response, err error) {
	timer := stopwatch.Start("SendMessage", stopwatch.LogStop)
	defer func() {
		timer.StopE(err)
	}()

	err = message.Validate()
	if err != nil {
		return nil, err
	}

	msMessage := ms.client.Email.NewMessage()

	from := mailersend.From{
		Name:  message.FromEmail.Name,
		Email: message.FromEmail.Address,
	}

	recipientList := []mailersend.Recipient{}
	for _, recipient := range message.Recipients {
		recipientList = append(recipientList, mailersend.Recipient{
			Email: recipient.Address,
			Name:  recipient.Name,
		})
	}

	for _, attachment := range message.Attachments {
		attachment := mailersend.Attachment{Filename: attachment.Filename, Content: attachment.Base64Content}
		msMessage.AddAttachment(attachment)

	}

	msMessage.SetFrom(from)
	msMessage.SetRecipients(recipientList)
	msMessage.SetSubject(message.Subject)
	msMessage.SetHTML(message.HtmlContent)
	msMessage.SetText(message.PlainTextContent)

	return ms.post(msMessage)
}

func (ms *MailerSend) post(messasge *mailersend.Message) (response *Response, err error) {
	if ms.client == nil {
		ms.client = mailersend.NewMailersend(ms.token)
		if ms.client == nil {
			return nil, fmt.Errorf("Failed to create MailerSend client")
		}
		ms.logf("Created new MailerSend client")
	}

	ctx := context.Background()
	msResponse, err := ms.client.Email.Send(ctx, messasge)
	if err != nil {
		return nil, err
	}

	if msResponse == nil {
		return nil, fmt.Errorf("No results returned from MailerSend")
	}

	bodyStr, err := readBody(msResponse.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	mapString := mapStringSliceToString(msResponse.Header)
	ms.logf("Send email: status_code=%d, body=%s, headers=%s", msResponse.StatusCode, bodyStr, mapString)

	response = &Response{
		StatusCode: msResponse.StatusCode,
		Body:       bodyStr,
		Headers:    msResponse.Header,
	}

	if msResponse.StatusCode < 200 || msResponse.StatusCode >= 300 {
		return response, fmt.Errorf("Failed to send email: status_code=%d, body=%s, headers=%s", msResponse.StatusCode, bodyStr, mapString)
	}

	return response, nil
}

// logf logs message either via defined user logger or via system one if no user logger is defined.
func (ms *MailerSend) logf(f string, args ...interface{}) {
	if ms.Logger != nil {
		ms.Logger(f, args...)
	} else {
		log.Printf(f, args...)
	}
}
