package utils

import (
	"bytes"
	"crypto/tls"
	"html/template"
	"log"

	"andalalin/initializers"

	gomail "gopkg.in/gomail.v2"
)

type Verification struct {
	Code    string
	Name    string
	Subject string
}

type ResetPassword struct {
	Code    string
	Subject string
}

func SendEmailVerification(email string, data *Verification) {
	config, err := initializers.LoadConfig()

	if err != nil {
		log.Fatal("could not load config", err)
	}

	// Sender data.
	from := config.EmailFrom
	smtpPass := config.SMTPPass
	smtpUser := config.SMTPUser
	to := email
	smtpHost := config.SMTPHost
	smtpPort := config.SMTPPort

	t, err := template.ParseFiles("templates/verificationCode.html")
	if err != nil {
		log.Fatal("Error reading the email template:", err)
		return
	}

	buffer := new(bytes.Buffer)
	if err = t.Execute(buffer, data); err != nil {
		log.Fatal("Error reading the email template:", err)
		return
	}

	m := gomail.NewMessage()

	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", data.Subject)
	m.SetBody("text/html", buffer.String())

	m.Embed("assets/andalalin.png")

	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send Email
	if err := d.DialAndSend(m); err != nil {
		log.Fatal("Could not send email: ", err)
	}
}

func SendEmailReset(email string, data *ResetPassword) {
	config, err := initializers.LoadConfig()

	if err != nil {
		log.Fatal("could not load config", err)
	}

	// Sender data.
	from := config.EmailFrom
	smtpPass := config.SMTPPass
	smtpUser := config.SMTPUser
	to := email
	smtpHost := config.SMTPHost
	smtpPort := config.SMTPPort

	t, err := template.ParseFiles("templates/resetCode.html")
	if err != nil {
		log.Fatal("Error reading the email template:", err)
		return
	}

	buffer := new(bytes.Buffer)
	if err = t.Execute(buffer, data); err != nil {
		log.Fatal("Error reading the email template:", err)
		return
	}

	m := gomail.NewMessage()

	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", data.Subject)
	m.SetBody("text/html", buffer.String())

	m.Embed("assets/andalalin.png")

	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send Email
	if err := d.DialAndSend(m); err != nil {
		log.Fatal("Could not send email: ", err)
	}
}
