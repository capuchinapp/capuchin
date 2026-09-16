package mailer

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"html/template"

	"github.com/wneessen/go-mail"
)

//go:embed logo.png templates/*
var efs embed.FS

// Mailer.
type Mailer struct {
	from   string
	client MailClient
}

// New Mailer.
func New(from string, client MailClient) *Mailer {
	return &Mailer{
		from:   from,
		client: client,
	}
}

// Отправить письмо.
func (m *Mailer) Send(ctx context.Context, recipients []string, subject string, tpl Template, data any) error {
	htpl, err := getTemplate(tpl)
	if err != nil {
		return fmt.Errorf("get template: %v", err)
	}

	content, err := assembleContent(htpl, data)
	if err != nil {
		return fmt.Errorf("assemble content: %v", err)
	}

	msg, err := assembleMessage(m.from, recipients, subject, content)
	if err != nil {
		return fmt.Errorf("assemble message: %v", err)
	}

	if err := m.client.DialAndSendWithContext(ctx, msg); err != nil {
		return fmt.Errorf("send mail: %v", err)
	}

	return nil
}

// getTemplate возвращает HTML-шаблон.
func getTemplate(tpl Template) (*template.Template, error) {
	switch tpl { //nolint:revive // ignore
	case TemplateAuthCode:
	case TemplateRegisterSuccess:
		// ok

	default:
		return nil, fmt.Errorf("unknown template: %s", tpl)
	}

	htpl, err := template.New("html_template").ParseFS(efs, LayoutBase.String(), tpl.String())
	if err != nil {
		return nil, fmt.Errorf("failed to parse html template: %s", err)
	}

	return htpl, nil
}

// assembleContent собирает содержимое письма.
func assembleContent(htpl *template.Template, data any) (string, error) {
	var buf bytes.Buffer

	err := htpl.ExecuteTemplate(&buf, LayoutBaseFilename, data)
	if err != nil {
		return "", fmt.Errorf("execute template: %v", err)
	}

	return buf.String(), nil
}

func assembleMessage(
	from string,
	recipients []string,
	subject string,
	htmlContent string,
) (*mail.Msg, error) {
	msg := mail.NewMsg()

	if err := msg.From(from); err != nil {
		return nil, fmt.Errorf("failed to set From address: %s", err)
	}

	if err := msg.To(recipients...); err != nil {
		return nil, fmt.Errorf("failed to set To address: %s", err)
	}

	if err := msg.EmbedFromEmbedFS("logo.png", &efs, mail.WithFileContentID("logo.png"), mail.WithFileContentType("image/png")); err != nil {
		return nil, fmt.Errorf("failed to embed logo: %s", err)
	}

	msg.Subject(subject)
	msg.SetBodyString(mail.TypeTextHTML, htmlContent)

	return msg, nil
}
