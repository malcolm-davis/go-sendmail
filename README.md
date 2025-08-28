# go-sendmail

## Overview

go-sendmail is a wrapper to mail service libraries.

### Features

- **sendmail interface**: Common interface to all mailclients
- **Cache**: Reuses the client instance across the application
- **MessageBuilder**: Provide a builder pattern to build messasges
- **stopwatch**: Auto stopwatch each message posted
- **Log override**: Provides logging override.  The default logging is slog

### Wrappers provided for

- SendGrid's Go Library at https://github.com/sendgrid/sendgrid-go
- MailerSend Go Lib at  https://github.com/mailersend/mailersend-go
- MailJet Go Lib at https://github.com/mailjet/mailjet-apiv3-go
- smtp2go Go Lib at https://github.com/smtp2go-oss/smtp2go-go
- mailtrap - wrapper around mailtrap json api request


### Install

```
go get github.com/malcolm-davis/go-sendmail
```

### Usage

#### Basic

```go
package main

import (
    "log/slog"
    "os"
    "time"

    "github.com/malcolm-davis/go-sendmail"
)

func main() {
    // Create a new stopwatch and start it
    send, err := sendmail.NewSendGrid(os.Getenv("SENDGRID_API_KEY"))
    if(err!=nil) {
        slog.Error("Error connecting to sendgrid service", err)
    }

    respnse, err := send.SendMail("From User", "test@example.com" "To User", "to@example.com")  "Sending with Twilio SendGrid is Fun",
         "and easy to do anywhere, even with Go", "<strong>and easy to do anywhere, even with Go</strong>") 
    if(err!=nil) {
        slog.Error("Error sending message", err)
    }
}
```

#### Message Builder

```go
package main

import (
    "log/slog"
    "os"
    "time"

    "github.com/malcolm-davis/go-sendmail"
)

func main() {
    fmt.Println("Starting mailjet example")
    fmt.Println("Make sure to set MAILJET_APIKEY and MAILJET_SECRETKEY environment variables")
    send, err := sendmail.NewMailJetMailManager(os.Getenv("MAILJET_API_KEY"), os.Getenv("MAILJET_SECRET_KEY"))
    if err != nil {
        slog.Error("Error connecting to mailjet service", "error", err)
        return
    }

    messageBuilder := sendmail.NewEmailMessage()

    // Note: the from email address requires authentication in MailJet
    messageBuilder.FromEmail("do-not-reply", "do-not-reply@example.dev")
    messageBuilder.AddRecipient("name", "user@example.com")
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
```


#### smtp2go

```
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
    send, err := sendmail.NewSmtp2go()
    if err != nil {
        slog.Error("Error connecting to smtp2go service", "error", err)
        return
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
```

#### Message Builder

```go
package main

import (
    "log/slog"
    "os"
    "time"

    "github.com/malcolm-davis/go-sendmail"
)

func main() {
    fmt.Println("Starting mailjet example")
    fmt.Println("Make sure to set MAILJET_APIKEY and MAILJET_SECRETKEY environment variables")
    send, err := sendmail.NewMailJet(os.Getenv("MAILJET_API_KEY"), os.Getenv("MAILJET_SECRET_KEY"))
    if err != nil {
        slog.Error("Error connecting to mailjet service", "error", err)
        return
    }

    messageBuilder := sendmail.NewEmailMessage()

    // Note: the from email address requires authentication in MailJet
    messageBuilder.FromEmail("do-not-reply", "do-not-reply@example.dev")
    messageBuilder.AddRecipient("name", "user@example.com")
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
```


#### MailTrap with image and multiple recipients

```

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
```
