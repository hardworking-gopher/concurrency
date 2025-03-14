package main

import (
	"bytes"
	"fmt"
	"github.com/vanng822/go-premailer/premailer"
	mail "github.com/xhit/go-simple-mail"
	"html/template"
	"sync"
	"time"
)

type Mailer struct {
	Domain      string
	Host        string
	Port        int
	Username    string
	Password    string
	Encryption  string
	FromAddress string
	FromName    string
	Wait        *sync.WaitGroup
	MailerChan  chan Message
	ErrorChan   chan error
	DoneChan    chan bool
}

type Message struct {
	From        string
	FromName    string
	To          string
	Subject     string
	Attachments []string
	Data        any
	DataMap     map[string]any
	Template    string
}

func (m *Mailer) sendMail(msg Message, errorChan chan error) {
	defer m.Wait.Done()

	if msg.Template == "" {
		msg.Template = "mail"
	}

	if msg.From == "" {
		msg.From = m.FromAddress
	}

	if msg.FromName == "" {
		msg.FromName = m.FromName
	}

	msg.DataMap = map[string]any{
		"message": msg.Data,
	}

	formattedMsg, err := m.buildHTMLMessage(msg)
	if err != nil {
		errorChan <- err
	}

	plainMsg, err := m.buildPlainMessage(msg)
	if err != nil {
		errorChan <- err
	}

	server := mail.NewSMTPClient()
	server.Host = m.Host
	server.Port = m.Port
	server.Username = m.Username
	server.Password = m.Password
	server.KeepAlive = false
	server.SendTimeout = 10 * time.Second
	server.ConnectTimeout = 10 * time.Second

	switch m.Encryption {
	case "tls":
		server.Encryption = mail.EncryptionTLS
	case "ssl":
		server.Encryption = mail.EncryptionSSL
	case "none":
		server.Encryption = mail.EncryptionNone
	}

	smtpClient, err := server.Connect()
	if err != nil {
		errorChan <- err
	}
	defer smtpClient.Close()

	email := mail.NewMSG()
	email.SetFrom(msg.From).AddTo(msg.To).SetSubject(msg.Subject)

	email.SetBody(mail.TextPlain, plainMsg)
	email.AddAlternative(mail.TextHTML, formattedMsg)

	if len(msg.Attachments) > 0 {
		for _, a := range msg.Attachments {
			email.AddAttachment(a)
		}
	}

	if err = email.Send(smtpClient); err != nil {
		errorChan <- err
	}
}

func (m *Mailer) buildHTMLMessage(msg Message) (string, error) {
	templateToRender := fmt.Sprintf("subscription_service/templates/%s.html.gohtml", msg.Template)

	t, err := template.New("email-html").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}

	formattedMsg, err := m.inlineCSS(tpl.String())
	if err != nil {
		return "", err
	}

	return formattedMsg, nil
}

func (m *Mailer) buildPlainMessage(msg Message) (string, error) {
	templateToRender := fmt.Sprintf("subscription_service/templates/%s.plain.gohtml", msg.Template)

	t, err := template.New("email-plain").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}

	return tpl.String(), nil
}

func (m *Mailer) inlineCSS(s string) (string, error) {
	options := &premailer.Options{
		RemoveClasses:     false,
		CssToAttributes:   false,
		KeepBangImportant: true,
	}

	p, err := premailer.NewPremailerFromString(s, options)
	if err != nil {
		return "", err
	}

	html, err := p.Transform()
	if err != nil {
		return "", err
	}

	return html, nil
}
