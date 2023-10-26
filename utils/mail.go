package utils

import (
	"bytes"
	"crypto/tls"
	"html/template"
	"log"

	"github.com/Ikit777/E-Andalalin/initializers"
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

type PersyaratanTidakSesuai struct {
	Jenis       string
	Kode        string
	Nama        string
	Tlp         string
	Status      string
	Persyaratan string
	Subject     string
}

type PermohonanSelesai struct {
	Kode    string
	Nama    string
	Tlp     string
	Jenis   string
	Status  string
	Subject string
}

type PermohonanDitolak struct {
	Kode    string
	Nama    string
	Tlp     string
	Jenis   string
	Status  string
	Subject string
}

type PermohonanDibatalkan struct {
	Kode    string
	Nama    string
	Tlp     string
	Jenis   string
	Status  string
	Subject string
}

type Pemasangan struct {
	Kode    string
	Nama    string
	Tlp     string
	Jenis   string
	Status  string
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

func SendEmailPersyaratan(email string, data *PersyaratanTidakSesuai) {
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

	t, err := template.ParseFiles("templates/persyaratanTidakSesuai.html")
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

func SendEmailPermohonanSelesai(email string, data *PermohonanSelesai) {
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

	t, err := template.ParseFiles("templates/permohonanSelesai.html")
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

func SendEmailPermohonanDibatalkan(email string, data *PermohonanDibatalkan) {
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

	t, err := template.ParseFiles("templates/permohonanDibatalkan.html")
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

func SendEmailPemasangan(email string, data *Pemasangan) {
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

	t, err := template.ParseFiles("templates/pemasanganPerlalin.html")
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

func SendEmailPermohonanDitolak(email string, data *PermohonanDitolak) {
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

	t, err := template.ParseFiles("templates/permohonanDitolak.html")
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
