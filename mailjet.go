// Wrapper to MailJet Go Library
// https://github.com/mailjet/mailjet-apiv3-go
package sendmail

import (
	"fmt"
	"log"

	"github.com/mailjet/mailjet-apiv3-go/v4"
	"github.com/malcolm-davis/go-stopwatch"
)

// @ref https://github.com/mailjet/mailjet-apiv3-go

// SendMailConfig holds the API key for MailJet
type MailJetConfig struct {
	APIKey    string
	SecretKey string
}

// SendMailManager provides functionality to send emails via MailJet
type MailJetMailManager struct {
	client    *mailjet.Client
	APIKey    string
	SecretKey string

	// User defined logger function.
	Logger func(string, ...interface{})
}

func NewMailJet(apiKey string, secretKey string) (*MailJetMailManager, error) {
	client := mailjet.NewMailjetClient(apiKey, secretKey)
	if client == nil {
		return nil, fmt.Errorf("Failed to create Mailjet client")
	}

	manager := &MailJetMailManager{
		APIKey:    apiKey,
		SecretKey: secretKey,
		client:    client,
	}

	return manager, nil
}

func (mj *MailJetMailManager) SendMail(fromName, fromEmail, toName, toEmail, subject, plainTextContent, htmlContent string) (response *Response, err error) {
	timer := stopwatch.Start("SendMail", stopwatch.LogStop)
	defer func() {
		timer.StopE(err)
	}()

	messagesInfo := []mailjet.InfoMessagesV31{{
		From: &mailjet.RecipientV31{
			Email: fromEmail,
			Name:  fromName,
		},
		To: &mailjet.RecipientsV31{
			mailjet.RecipientV31{
				Email: toEmail,
				Name:  toName,
			},
		},
		Subject:  subject,
		TextPart: plainTextContent,
		HTMLPart: htmlContent,
	},
	}

	return mj.post(messagesInfo)
}

func (mj *MailJetMailManager) SendMessage(message *Message) (response *Response, err error) {
	timer := stopwatch.Start("SendMessage", stopwatch.LogStop)
	defer func() {
		timer.StopE(err)
	}()

	if len(message.Recipients) == 0 {
		return nil, fmt.Errorf("missing recipient(s) address")
	}

	if message.FromEmail == nil {
		return nil, fmt.Errorf("missing from email address")
	}

	recipientList := mailjet.RecipientsV31{}
	for _, recipient := range message.Recipients {
		recipientList = append(recipientList, mailjet.RecipientV31{
			Email: recipient.Address,
			Name:  recipient.Name,
		})
	}

	attachmentList := mailjet.AttachmentsV31{}
	for _, attachment := range message.Attachments {
		attachmentList = append(attachmentList, mailjet.AttachmentV31{
			ContentType:   attachment.ContentType,
			Filename:      attachment.Filename,
			Base64Content: attachment.Base64Content,
		})
	}

	messagesInfo := []mailjet.InfoMessagesV31{{
		From: &mailjet.RecipientV31{
			Email: message.FromEmail.Address,
			Name:  message.FromEmail.Name,
		},
		To:          &recipientList,
		Subject:     message.Subject,
		TextPart:    message.PlainTextContent,
		HTMLPart:    message.HtmlContent,
		Attachments: &attachmentList,
	},
	}

	return mj.post(messagesInfo)
}

func (mj *MailJetMailManager) post(messagesInfo []mailjet.InfoMessagesV31) (response *Response, err error) {
	client := mj.client
	if client == nil {
		client := mailjet.NewMailjetClient(mj.APIKey, mj.SecretKey)
		if client == nil {
			return nil, fmt.Errorf("Failed to create Mailjet client")
		}
		mj.logf("Created new Mailjet client")
	}

	messages := mailjet.MessagesV31{Info: messagesInfo}
	mailjetResponse, err := client.SendMailV31(&messages)
	if err != nil {
		return nil, err
	}

	if mailjetResponse == nil || len(mailjetResponse.ResultsV31) == 0 {
		return nil, fmt.Errorf("No results returned from Mailjet")
	}

	statusCode := 400
	if mailjetResponse.ResultsV31[0].Status == "success" {
		statusCode = 200
	}

	if statusCode == 200 {
		mj.logf("Send email: status_code=%d, result=%v", statusCode, mailjetResponse.ResultsV31[0])
	}

	response = &Response{
		StatusCode: statusCode,
	}
	// Accept any 2xx response as success
	if statusCode < 200 || statusCode >= 300 {
		return response, fmt.Errorf("Failed to send email: %s", response.Body)
	}
	return response, nil
}

// logf logs message either via defined user logger or via system one if no user logger is defined.
func (mj *MailJetMailManager) logf(f string, args ...interface{}) {
	if mj.Logger != nil {
		mj.Logger(f, args...)
	} else {
		log.Printf(f, args...)
	}
}
