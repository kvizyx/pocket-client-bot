package telegram

import (
	"context"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/kviz48/pocket-client-bot/internal/storage"
	"github.com/zhashkevych/go-pocket-sdk"
	"net/url"
)

const (
	Start = "start"
)

func (b *Bot) handeCommand(msg *tgbotapi.Message) error {
	switch msg.Command() {
	case Start:
		return b.handleStartCommand(msg)
	default:
		return b.handleUnknownCommand(msg)
	}
}

func (b *Bot) handeMessage(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, b.messagesConfig.Responses.LinkSaveSuccess)

	_, err := url.ParseRequestURI(message.Text)

	if err != nil {
		return errInvalidURL
	}

	accessToken, err := b.storage.Get(message.Chat.ID, storage.AccessTokens)
	if err != nil {
		return errUnauthorized
	}

	err = b.pocket.Add(context.Background(), pocket.AddInput{
		URL:         message.Text,
		AccessToken: accessToken,
	})
	if err != nil {
		return errCantSave
	}

	_, err = b.tg.Send(msg)
	return err
}

func (b *Bot) handleStartCommand(message *tgbotapi.Message) error {
	_, err := b.getAccessToken(message.Chat.ID)

	if err != nil {
		return b.authorizeNewUser(message)
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, b.messagesConfig.Responses.AlreadyAuthorized)
	_, err = b.tg.Send(msg)

	return err
}

func (b *Bot) handleUnknownCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, b.messagesConfig.Responses.UnknownCommand)

	_, err := b.tg.Send(msg)
	return err
}
