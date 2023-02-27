package telegram

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/kviz48/pocket-client-bot/internal/storage"
	"log"
)

func (b *Bot) authorizeNewUser(message *tgbotapi.Message) error {
	authLink, err := b.generateAuthURL(message.Chat.ID)
	if err != nil {
		return err
	}

	reply := fmt.Sprintf(b.messagesConfig.Responses.Start, message.From.UserName, authLink)
	msg := tgbotapi.NewMessage(message.Chat.ID, reply)

	_, err = b.tg.Send(msg)
	return err
}

func (b *Bot) generateAuthURL(chatID int64) (string, error) {
	redirectURL := b.generateRedirectURL(chatID)

	reqToken, err := b.pocket.GetRequestToken(context.Background(), redirectURL)
	if err != nil {
		return "", err
	}

	err = b.storage.Save(chatID, reqToken, storage.RequestTokens)
	if err != nil {
		log.Fatal("cannot save request token")
	}

	authLink, err := b.pocket.GetAuthorizationURL(reqToken, redirectURL)
	if err != nil {
		return "", err
	}

	return authLink, nil
}

func (b *Bot) generateRedirectURL(chatID int64) (url string) {
	url = fmt.Sprintf("%s?chat_id=%d", b.redirectURL, chatID)
	return
}

func (b *Bot) getAccessToken(chatID int64) (string, error) {
	return b.storage.Get(chatID, storage.AccessTokens)
}
