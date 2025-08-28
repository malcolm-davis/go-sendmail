package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/malcolm-davis/go-sendmail"
)

func main() {
	fmt.Println("Starting mailjet example")
	fmt.Println("Make sure to set MAILJET_API_KEY and MAILJET_SECRET_KEY environment variables")
	send, err := sendmail.NewMailJet(os.Getenv("MAILJET_API_KEY"), os.Getenv("MAILJET_SECRET_KEY"))
	if err != nil {
		slog.Error("Error connecting to mailjet service", "error", err)
		return
	}

	messageBuilder := sendmail.NewEmailMessage()

	// Note: the from email address requires authentication in MailJet
	messageBuilder.FromEmail("do-not-reply", "do-not-reply@example.com")
	messageBuilder.AddRecipient("user", "user@example.com")
	messageBuilder.AddRecipient("name 2", "user2@example.com")
	messageBuilder.Subject("contact us email")
	messageBuilder.PlainTextContent("Test mailjet")
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
