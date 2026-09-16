package mailer

import (
	"context"

	"github.com/wneessen/go-mail"
)

// MailClient почтовый клиент.
type MailClient interface {
	// Отправить письмо.
	DialAndSendWithContext(ctx context.Context, messages ...*mail.Msg) error
}
