// Wrapper to Smtp2go Go Library
// https://github.com/Smtp2go/Smtp2go-apiv3-go
package sendmail

import (
	"fmt"
	"log"

	"github.com/malcolm-davis/go-stopwatch"
	"github.com/smtp2go-oss/smtp2go-go"
)

// @ref https://github.com/Smtp2go/Smtp2go-apiv3-go

// SendMailConfig holds the API key for Smtp2go

// SendMailManager provides functionality to send emails via Smtp2go
type Smtp2goMail struct {
	APIToken string

	// User defined logger function.
	Logger func(string, ...interface{})
}

func NewSmtp2go(smtp2goKey string) (*Smtp2goMail, error) {
	manager := &Smtp2goMail{
		APIToken: smtp2goKey,
	}

	return manager, nil
}

func (ms *Smtp2goMail) SendMail(fromName, fromEmail, toName, toEmail, subject, plainTextContent, htmlContent string) (response *Response, err error) {
	timer := stopwatch.Start("SendMail", stopwatch.LogStop)
	defer func() {
		timer.StopE(err)
	}()

	from := fmt.Sprintf("%s <%s>", fromName, fromEmail)
	message := smtp2go.Email{
		From:     from,
		To:       []string{fmt.Sprintf("%s <%s>", toName, toEmail)},
		Subject:  subject,
		TextBody: plainTextContent,
		HtmlBody: htmlContent,
	}

	return ms.post(message)
}

func (ms *Smtp2goMail) SendMessage(message *Message) (response *Response, err error) {
	timer := stopwatch.Start("SendMessage", stopwatch.LogStop)
	defer func() {
		timer.StopE(err)
	}()

	err = message.Validate()
	if err != nil {
		return nil, err
	}

	from := fmt.Sprintf("%s <%s>", message.FromEmail.Name, message.FromEmail.Address)

	recipientList := []string{}
	for _, recipient := range message.Recipients {
		recipientList = append(recipientList, fmt.Sprintf("%s <%s>", recipient.Name, recipient.Address))
	}

	attachmentList := []*smtp2go.EmailBinaryData{}
	for _, attachment := range message.Attachments {
		attachmentList = append(attachmentList, &smtp2go.EmailBinaryData{
			Filename: attachment.Filename,
			Fileblob: attachment.Base64Content,
			MimeType: attachment.ContentType,
		})
	}

	email := smtp2go.Email{
		From:        from,
		To:          recipientList,
		Subject:     message.Subject,
		TextBody:    message.PlainTextContent,
		HtmlBody:    message.HtmlContent,
		Attachments: attachmentList,
	}

	return ms.post(email)
}

func (ms *Smtp2goMail) post(email smtp2go.Email) (response *Response, err error) {

	res, err := smtp2go.Send(&email)
	if err != nil {
		return nil, err
	}

	if res == nil {
		return nil, fmt.Errorf("No results returned from Smtp2go")
	}

	if len(res.Data.Error) != 0 {
		return nil, fmt.Errorf("Smtp2go error: %s", res.Data.Error)
	}

	response = &Response{
		StatusCode: 200,
		Body:       fmt.Sprintf("RequestId: %s", res.RequestId),
	}

	return response, nil
}

// logf logs message either via defined user logger or via system one if no user logger is defined.
func (ms *Smtp2goMail) logf(f string, args ...interface{}) {
	if ms.Logger != nil {
		ms.Logger(f, args...)
	} else {
		log.Printf(f, args...)
	}
}
