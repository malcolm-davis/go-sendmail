package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/malcolm-davis/go-sendmail"
)

func main() {
	fmt.Println("Make sure you have set the SENDGRID_API_KEY environment variable")

	// Create a new sendmail manager and init
	send, err := sendmail.NewSendGrid(os.Getenv("SENDGRID_API_KEY"))
	if err != nil {
		slog.Error("Error connecting to sendgrid service", "error", err)
		return
	}

	response, err := send.SendMail("do-not-reply", "do-not-reply@example.com", "touser", "user@example.com", " contact us email",
		"Test go-sendemail", "<strong>and easy to do anywhere, even with Go</strong>")
	if err != nil {
		slog.Error("Error sending message", "error", err)
		return
	}

	fmt.Println("response.StatusCode", response.StatusCode)

}
