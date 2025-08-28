package sendmail

import (
	"io"
	"strings"
)

// Response holds the response from an API call.
type Response struct {
	StatusCode int                 // e.g. 200
	Body       string              // e.g. {"result: success"}
	Headers    map[string][]string // e.g. map[X-Ratelimit-Limit:[600]]
}

type SendMail interface {
	// send mail is intended for a single recipient
	SendMail(fromName, fromEmail, toName, toEmail, subject, plainTextContent, htmlContent string) (response *Response, err error)

	// SendMessage allows for more complex email scenarios, including sending emails to multiple recipients and attachments
	SendMessage(message *Message) (response *Response, err error)
}

func readBody(body io.ReadCloser) (string, error) {
	defer body.Close()
	buf, err := io.ReadAll(body)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

func mapStringSliceToString(m map[string][]string) string {
	var sb strings.Builder

	for key, values := range m {
		sb.WriteString(key)
		sb.WriteString(": [")
		for i, value := range values {
			sb.WriteString("\"")
			sb.WriteString(value)
			sb.WriteString("\"")
			if i < len(values)-1 {
				sb.WriteString(", ")
			}
		}
		sb.WriteString("]")
		sb.WriteString(" ")
	}

	return sb.String()
}
