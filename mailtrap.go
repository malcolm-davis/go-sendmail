// Wrapper to Smtp2go Go Library
// https://github.com/Smtp2go/Smtp2go-apiv3-go
package sendmail

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"net/http"
	"time"

	"github.com/malcolm-davis/go-stopwatch"
)

// MailTrapConfig provides functionality to send emails via mailtrap
// https://mailtrap.io/blog/golang-send-email/#Send-emails-in-Go-using-email-API
type MailTrap struct {
	token  string
	client *http.Client
}

func NewMailTrap(mailTrapKey string) (*MailTrap, error) {
	client := http.Client{Timeout: 10 * time.Second}
	manager := &MailTrap{
		token:  mailTrapKey,
		client: &client,
	}

	return manager, nil
}

func (ms *MailTrap) SendMail(fromName, fromEmail, toName, toEmail, subject, plainTextContent, htmlContent string) (response *Response, err error) {
	timer := stopwatch.Start("SendMail", stopwatch.LogStop)
	defer func() {
		timer.StopE(err)
	}()

	// format
	// message := []byte(`{
	//     "from":{"email":"john.doe@your.domain"},
	//     "to":[{"email":"kate.doe@example.com"}],
	//     "subject":"Why aren’t you using Mailtrap yet?",
	//     "text":"Here’s the space for your great sales pitch",
	//     "html":"<strong>Here’s the space for your great sales pitch</strong>"
	// }`)

	// from := fmt.Sprintf(`"from":{"email":"%s"}`, fromEmail)
	// to := fmt.Sprintf(`"to":[{"email":"%s"}]`, toEmail)
	// sub := fmt.Sprintf(`"subject":"%s"`, subject)
	// content := fmt.Sprintf(`"text":"%s"`, plainTextContent)
	// html := fmt.Sprintf(`{"html":%s"}`, htmlContent)
	// message := []byte(fmt.Sprintf("{%s,%s,%s,%s,%s}", from, to, sub, content, html))

	messageBuilder := NewEmailMessage()
	messageBuilder.FromEmail(fromName, fromEmail)
	messageBuilder.AddRecipient(toName, toEmail)
	messageBuilder.Subject(subject)
	messageBuilder.PlainTextContent(plainTextContent)
	messageBuilder.HtmlContent(htmlContent)

	message, err := messageBuilder.Build()
	if err != nil {
		return nil, err
	}

	email, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}
	return ms.post(email)
}

func (ms *MailTrap) SendMessage(message *Message) (response *Response, err error) {
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

	//  message := []byte(`{
	//     "from":{"email":"john.doe@your.domain"},
	//     "to":[
	//         {"email":"kate.doe@example.com"},
	//         {"email":"alex.doe@example.com"},
	//         {"email":"lisa.doe@example.com"}
	//     ],
	//     "subject":"Why aren’t you using Mailtrap yet?",
	//     "text":"Here’s the space for your great sales pitch",
	//     "html":"<strong>Here’s the space for your great sales pitch</strong>"
	// }`)

	// message := []byte(`{
	//     "from": { "email": "john.doe@your.domain" },
	//     "to": [
	//         { "email": "kate.doe@example.com" }
	//     ],
	//     "subject": "Here’s your attached file!",
	//     "text": "Check out the attached file.",
	//     "html": "<p>Check out the attached <strong>file</strong>.</p>",
	//     "attachments": [
	//         {
	//           "filename": "example.pdf",
	//           "content": "` + encodedFileData + `",
	//           "type": "application/pdf",
	//           "disposition": "attachment"
	//         }
	//     ]
	// }`)

	email, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}
	return ms.post(email)
}

func (ms *MailTrap) post(message []byte) (response *Response, err error) {
	httpHost := "https://send.api.mailtrap.io/api/send"
	request, err := http.NewRequest(http.MethodPost, httpHost, bytes.NewBuffer(message))
	if err != nil {
		return nil, err
	}
	// Set required headers
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+ms.token)

	// Send request
	if ms.client == nil {
		ms.client = &http.Client{Timeout: 10 * time.Second}
	}

	res, err := ms.client.Do(request)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, fmt.Errorf("No results returned from mailtrap.io")
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Mailtrap API error: status code %d, body: %s", res.StatusCode, string(body))
	}

	response = &Response{
		StatusCode: res.StatusCode,
		Body:       string(body),
		Headers:    res.Header,
	}

	return response, nil
}
