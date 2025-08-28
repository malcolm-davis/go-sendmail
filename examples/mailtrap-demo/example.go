package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/malcolm-davis/go-sendmail"
)

func main() {
	fmt.Println("Make sure you have set the MAILTRAP_API_KEY in an environment variable")

	// Create a new sendmail manager and init
	send, err := sendmail.NewMailTrap(os.Getenv("MAILTRAP_API_KEY"))
	if err != nil {
		slog.Error("Error connecting to sendgrid service", "error", err)
		return
	}

	// Message builder example
	messageBuilder := sendmail.NewEmailMessage()
	messageBuilder.FromEmail("do-not-reply", "do-not-reply@example.dev")
	messageBuilder.AddRecipient("Use ", "user_1@example.com")
	messageBuilder.AddRecipient("User 2", "user_2@example.com")
	messageBuilder.Subject("Transactional email from mailtrap")
	messageBuilder.PlainTextContent("Test go-sendemail mailtrap feature")
	messageBuilder.AddAttachment("image/png", "test.png",
		"iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAAABHNCSVQICAgIfAhkiAAAAD1JREFUOI1jYKAQMJKhJxKJvZyJUhewEKEG2ZJ/5BjwF4mN4WWKvUCTMPiPxCYYS4PDCyQ5meouGAYGUAwAmxIFIoTEH0QAAAAASUVORK5CYII=")

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
	fmt.Println("response.StatusCode", response.StatusCode)

}
