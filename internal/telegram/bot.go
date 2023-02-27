package telegram

import (
	"errors"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/kviz48/pocket-client-bot/internal/config"
	"github.com/kviz48/pocket-client-bot/internal/storage/boltdb"
	"github.com/zhashkevych/go-pocket-sdk"
	"log"
	"strings"
)

type Bot struct {
	tg             *tgbotapi.BotAPI
	pocket         *pocket.Client
	storage        *boltdb.TokenStorage
	messagesConfig *config.MessagesConfig
	redirectURL    string
}

func NewBot(tg *tgbotapi.BotAPI, pocket *pocket.Client, redirectURL string, ts *boltdb.TokenStorage, messagesCfg *config.MessagesConfig) (*Bot, error) {
	if strings.TrimSpace(redirectURL) == "" {
		return &Bot{}, errors.New("empty redirect url")
	}

	return &Bot{
		tg:             tg,
		pocket:         pocket,
		storage:        ts,
		messagesConfig: messagesCfg,
		redirectURL:    redirectURL,
	}, nil
}

func (b *Bot) Start() error {
	log.Printf("Authorized on account %s", b.tg.Self.UserName)

	updates, err := b.initUpdatesChannel()
	if err != nil {
		return err
	}

	b.handleUpdates(updates)

	return nil
}

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			if err := b.handeCommand(update.Message); err != nil {
				b.handleError(update.Message.Chat.ID, err)
			}
			continue
		}

		if err := b.handeMessage(update.Message); err != nil {
			b.handleError(update.Message.Chat.ID, err)
		}
	}
}

func (b *Bot) initUpdatesChannel() (tgbotapi.UpdatesChannel, error) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	return b.tg.GetUpdatesChan(u)
}
