package sendmail

import (
	"fmt"
	"strings"
)

// Email holds email name and address info
type Email struct {
	Name    string `json:"name,omitempty"`
	Address string `json:"email,omitempty"`
}

type Attachment struct {
	ContentType   string `json:"type,omitempty"`
	Filename      string `json:"filename,omitempty"`
	Base64Content string `json:"content,omitempty"`
	Disposition   string `json:"disposition,omitempty"`
}

type Message struct {
	FromEmail        *Email        `json:"from,omitempty"`
	Recipients       []*Email      `json:"to,omitempty"`
	Subject          string        `json:"subject,omitempty"`
	PlainTextContent string        `json:"text,omitempty"`
	HtmlContent      string        `json:"html,omitempty"`
	Attachments      []*Attachment `json:"attachments,omitempty"`
}

type MessageBuilder interface {
	FromEmail(name, address string) MessageBuilder
	AddRecipient(name, address string) MessageBuilder
	Subject(subject string) MessageBuilder
	PlainTextContent(plainTextContent string) MessageBuilder
	HtmlContent(htmlContent string) MessageBuilder
	AddAttachment(contentType, filename, base64Content string, disposition_optional ...string) MessageBuilder
	Build() (*Message, error)
}

type messageBuilder struct {
	emailMessage *Message
}

func NewEmailMessage() MessageBuilder {
	return &messageBuilder{emailMessage: &Message{}}
}

func (m *messageBuilder) FromEmail(name, address string) MessageBuilder {
	if strings.TrimSpace(address) == "" {
		return m
	}
	m.emailMessage.FromEmail = &Email{Name: name, Address: address}
	return m
}

func (m *messageBuilder) AddRecipient(name, address string) MessageBuilder {
	if strings.TrimSpace(address) == "" {
		return m
	}
	m.emailMessage.Recipients = append(m.emailMessage.Recipients, &Email{Name: name, Address: address})
	return m
}

func (m *messageBuilder) Subject(subject string) MessageBuilder {
	m.emailMessage.Subject = subject
	return m
}

func (m *messageBuilder) PlainTextContent(content string) MessageBuilder {
	m.emailMessage.PlainTextContent = content
	return m
}

func (m *messageBuilder) HtmlContent(content string) MessageBuilder {
	m.emailMessage.HtmlContent = content
	return m
}

func (m *messageBuilder) AddAttachment(contentType, filename, base64Content string, disposition_optional ...string) MessageBuilder {
	disposition := "attachment"
	if len(disposition_optional) > 0 {
		disposition = disposition_optional[0]
	}
	m.emailMessage.Attachments = append(m.emailMessage.Attachments, &Attachment{
		ContentType:   contentType,
		Filename:      filename,
		Base64Content: base64Content,
		Disposition:   disposition,
	})
	return m
}

func (m *messageBuilder) Build() (*Message, error) {
	if m.emailMessage.FromEmail == nil {
		return nil, fmt.Errorf("missing from email address")
	}
	if len(m.emailMessage.Recipients) == 0 {
		return nil, fmt.Errorf("missing recipient(s) address")
	}
	if m.emailMessage.Subject == "" {
		return nil, fmt.Errorf("missing subject")
	}
	return m.emailMessage, nil
}
