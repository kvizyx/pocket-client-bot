package telegram

import (
	"errors"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	errInvalidURL   = errors.New("invalid url")
	errUnauthorized = errors.New("you are not authorized")
	errCantSave     = errors.New("cant save link :(")
)

func (b *Bot) handleError(chatID int64, err error) {
	msg := tgbotapi.NewMessage(chatID, b.messagesConfig.Errors.Default)

	switch err {
	case errInvalidURL:
		msg.Text = b.messagesConfig.Errors.InvalidLink
		b.tg.Send(msg)
	case errUnauthorized:
		msg.Text = b.messagesConfig.Errors.Unauthorized
		b.tg.Send(msg)
	default:
		b.tg.Send(msg)
	}
}
