package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/malcolm-davis/go-sendmail"
)

func main() {
	fmt.Println("Starting smtp2go example")
	fmt.Println("Make sure to set SMTP2GO_API_KEY")
	send, err := sendmail.NewSmtp2go(os.Getenv("SMTP2GO_API_KEY"))
	if err != nil {
		slog.Error("Error connecting to smtp2go service", "error", err)
		return
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	send.Logger = func(format string, args ...interface{}) {
		// slog.Debug(fmt.Sprintf(format, args...))
		logger.Info(format, args...)
	}
	messageBuilder := sendmail.NewEmailMessage()

	// Note: the from email address requires authentication in Smtp2go
	messageBuilder.FromEmail("do-not-reply", "do-not-reply@example.com")
	messageBuilder.AddRecipient("user", "user@example.com")
	messageBuilder.Subject("contact us email")
	messageBuilder.PlainTextContent("Test smtp2go")
	messageBuilder.HtmlContent("<strong>and easy to do anywhere, even with Go</strong>")

	message, err := messageBuilder.Build()
	if err != nil {
		slog.Error("Error building message", "error", err)
		return
	}
	response, err := send.SendMessage(message)
	if err != nil {
		slog.Error("Error sending message", "error", err)
		return
	}

	// if the response is successful, print the status code 200
	// however, if the from email address is not a verified domain in mailjet, a 200 will be returned, but no mail delivery will occur
	fmt.Println("response.StatusCode", response.StatusCode)

}
