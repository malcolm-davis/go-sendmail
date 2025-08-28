package sendmail

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEmailMessage(t *testing.T) {
	builder := NewEmailMessage()
	assert.NotNil(t, builder)
}

func TestMessageBuilder_Build_Success(t *testing.T) {
	builder := NewEmailMessage()
	message, err := builder.
		FromEmail("Sender", "sender@example.com").
		AddRecipient("Recipient", "recipient@example.com").
		Subject("Test Subject").
		PlainTextContent("Plain text content").
		HtmlContent("<p>HTML content</p>").
		Build()

	require.NoError(t, err)
	assert.NotNil(t, message)
	assert.Equal(t, "Sender", message.FromEmail.Name)
	assert.Equal(t, "sender@example.com", message.FromEmail.Address)
	assert.Len(t, message.Recipients, 1)
	assert.Equal(t, "Recipient", message.Recipients[0].Name)
	assert.Equal(t, "recipient@example.com", message.Recipients[0].Address)
	assert.Equal(t, "Test Subject", message.Subject)
	assert.Equal(t, "Plain text content", message.PlainTextContent)
	assert.Equal(t, "<p>HTML content</p>", message.HtmlContent)
}

func TestMessageBuilder_Build_MissingFromEmail(t *testing.T) {
	builder := NewEmailMessage()
	_, err := builder.
		AddRecipient("Recipient", "recipient@example.com").
		Subject("Test Subject").
		Build()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing from email address")
}

func TestMessageBuilder_Build_MissingRecipient(t *testing.T) {
	builder := NewEmailMessage()
	_, err := builder.
		FromEmail("Sender", "sender@example.com").
		Subject("Test Subject").
		Build()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing recipient(s) address")
}

func TestMessageBuilder_Build_MissingSubject(t *testing.T) {
	builder := NewEmailMessage()
	_, err := builder.
		FromEmail("Sender", "sender@example.com").
		AddRecipient("Recipient", "recipient@example.com").
		Build()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing subject")
}

func TestMessageBuilder_EmptyEmailAddress(t *testing.T) {
	builder := NewEmailMessage()
	message, err := builder.
		FromEmail("Sender", "sender@example.com").
		AddRecipient("", ""). // Empty recipient should be ignored
		AddRecipient("Valid", "valid@example.com").
		Subject("Test Subject").
		Build()

	require.NoError(t, err)
	assert.Len(t, message.Recipients, 1)
	assert.Equal(t, "Valid", message.Recipients[0].Name)
}

func TestMessageBuilder_AddAttachment(t *testing.T) {
	builder := NewEmailMessage()
	message, err := builder.
		FromEmail("Sender", "sender@example.com").
		AddRecipient("Recipient", "recipient@example.com").
		Subject("Test Subject").
		AddAttachment("text/plain", "test.txt", "dGVzdCBjb250ZW50", "attachment").
		Build()

	require.NoError(t, err)
	assert.Len(t, message.Attachments, 1)
	assert.Equal(t, "text/plain", message.Attachments[0].ContentType)
	assert.Equal(t, "test.txt", message.Attachments[0].Filename)
	assert.Equal(t, "dGVzdCBjb250ZW50", message.Attachments[0].Base64Content)
	assert.Equal(t, "attachment", message.Attachments[0].Disposition)
}

func TestMessageBuilder_AddAttachment_DefaultDisposition(t *testing.T) {
	builder := NewEmailMessage()
	message, err := builder.
		FromEmail("Sender", "sender@example.com").
		AddRecipient("Recipient", "recipient@example.com").
		Subject("Test Subject").
		AddAttachment("text/plain", "test.txt", "dGVzdCBjb250ZW50").
		Build()

	require.NoError(t, err)
	assert.Len(t, message.Attachments, 1)
	assert.Equal(t, "attachment", message.Attachments[0].Disposition)
}
