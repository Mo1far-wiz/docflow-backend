package emailer

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"mime/quotedprintable"
	"net/smtp"
	"os"
	"path/filepath"
	"strings"
)

func SendEmail(to, subject, body, attachment string) error {
	from := os.Getenv("SMTP_USER")
	password := os.Getenv("SMTP_PASS")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	// Create a buffer to hold the email message
	var msg strings.Builder

	// Create a new multipart message
	multipartWriter := multipart.NewWriter(&msg)

	// Set the headers
	msg.WriteString(fmt.Sprintf("From: %s\r\n", from))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", to))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=\"%s\"\r\n", multipartWriter.Boundary()))
	msg.WriteString("\r\n")

	// Add the text part
	textPart, err := multipartWriter.CreatePart(
		map[string][]string{
			"Content-Type":              {"text/plain; charset=\"utf-8\""},
			"Content-Transfer-Encoding": {"quoted-printable"},
		},
	)
	if err != nil {
		return err
	}

	// Write the body into the text part
	_, err = quotedprintable.NewWriter(textPart).Write([]byte(body))
	if err != nil {
		return err
	}

	// Open the attachment file
	file, err := os.Open(attachment)
	if err != nil {
		return err
	}
	defer file.Close()

	// Get the file info to determine the file name and content type
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	fileName := filepath.Base(fileInfo.Name())

	// Add the attachment part
	attachmentPart, err := multipartWriter.CreatePart(
		map[string][]string{
			"Content-Type":              {"application/pdf"},
			"Content-Disposition":       {fmt.Sprintf("attachment; filename=\"%s\"", fileName)},
			"Content-Transfer-Encoding": {"base64"},
		},
	)
	if err != nil {
		return err
	}

	// Encode the attachment in base64
	attachmentData, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	base64Encoded := make([]byte, base64.StdEncoding.EncodedLen(len(attachmentData)))
	base64.StdEncoding.Encode(base64Encoded, attachmentData)

	// Write the base64 encoded data to the attachment part
	_, err = attachmentPart.Write(base64Encoded)
	if err != nil {
		return err
	}

	// Close the multipart writer
	multipartWriter.Close()

	// Set up the SMTP connection
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Create a custom TLS config
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true, // You might want to set this to false in production
		ServerName:         smtpHost,
	}

	// Connect to the SMTP server using TLS
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%s", smtpHost, smtpPort), tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to the SMTP server: %v", err)
	}

	// Create a new SMTP client
	client, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %v", err)
	}
	defer client.Quit()

	// Authenticate
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("authentication failed: %v", err)
	}

	// Set the sender and recipient
	if err = client.Mail(from); err != nil {
		return fmt.Errorf("failed to set sender: %v", err)
	}
	if err = client.Rcpt(to); err != nil {
		return fmt.Errorf("failed to set recipient: %v", err)
	}

	// Send the email
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to send email data: %v", err)
	}
	defer w.Close()

	_, err = w.Write([]byte(msg.String()))
	if err != nil {
		return fmt.Errorf("failed to write email data: %v", err)
	}

	return nil
}
