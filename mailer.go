package gomail

import (
	"io/ioutil"
	"net/smtp"
	"os"
	"path/filepath"
	"text/template"
)

type MailSender interface {
	SendMessage(messages ...*Message) (err error)
}

type Mailer struct {
	Server   string
	Port     int
	UserName string
	Password string
	Host     string    // This is optional, only used if you want to tell smtp server your hostname
	Auth     smtp.Auth // This is optional, only used if Authentication is not plain
	Sender   *Sender   // This is optional, only used if the From/ReplyTo does not specified in the message
}

type Sender struct {
	From    string
	ReplyTo string
}

var MailTemplate *template.Template

// Send the given email messages using this Mailer.
func (m *Mailer) SendMessage(messages ...*Message) (err error) {
	if m.Auth == nil {
		m.Auth = smtp.PlainAuth(m.UserName, m.UserName, m.Password, m.Server)
	}

	c, err := Transport(m.Server, m.Port, m.Host, m.Auth)
	if err != nil {
		return
	}
	defer c.Quit()

	for _, message := range messages {
		m.fillDefault(message)
		if err = Send(c, message); err != nil {
			return
		}
	}

	return
}

func TemplateFolder(templateFolder string) error {
	MailTemplate = template.New(templateFolder)
	err := filepath.Walk(templateFolder, func(path string, info os.FileInfo, err error) error {
		r, err := filepath.Rel(templateFolder, path)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			buf, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			tmpl := MailTemplate.New(filepath.ToSlash(r))

			template.Must(tmpl.Parse(string(buf)))
		}

		return nil
	})
	return err
}

func (m *Mailer) fillDefault(message *Message) {
	if m.Sender == nil {
		return
	}
	if message.From == "" {
		message.From = m.Sender.From
	}

	if message.ReplyTo == "" {
		message.ReplyTo = m.Sender.ReplyTo
	}
}
