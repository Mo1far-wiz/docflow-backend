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

	var msg strings.Builder

	multipartWriter := multipart.NewWriter(&msg)

	msg.WriteString(fmt.Sprintf("From: %s\r\n", from))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", to))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=\"%s\"\r\n", multipartWriter.Boundary()))
	msg.WriteString("\r\n")

	textPart, err := multipartWriter.CreatePart(
		map[string][]string{
			"Content-Type":              {"text/plain; charset=\"utf-8\""},
			"Content-Transfer-Encoding": {"quoted-printable"},
		},
	)
	if err != nil {
		return err
	}

	_, err = quotedprintable.NewWriter(textPart).Write([]byte(body))
	if err != nil {
		return err
	}

	file, err := os.Open(attachment)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	fileName := filepath.Base(fileInfo.Name())

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

	attachmentData, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	base64Encoded := make([]byte, base64.StdEncoding.EncodedLen(len(attachmentData)))
	base64.StdEncoding.Encode(base64Encoded, attachmentData)

	_, err = attachmentPart.Write(base64Encoded)
	if err != nil {
		return err
	}

	multipartWriter.Close()

	auth := smtp.PlainAuth("", from, password, smtpHost)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpHost,
	}

	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%s", smtpHost, smtpPort), tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to the SMTP server: %v", err)
	}

	client, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %v", err)
	}
	defer client.Quit()

	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("authentication failed: %v", err)
	}

	if err = client.Mail(from); err != nil {
		return fmt.Errorf("failed to set sender: %v", err)
	}
	if err = client.Rcpt(to); err != nil {
		return fmt.Errorf("failed to set recipient: %v", err)
	}

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
