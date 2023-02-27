package main

import (
	"github.com/boltdb/bolt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/kviz48/pocket-client-bot/internal/config"
	"github.com/kviz48/pocket-client-bot/internal/server"
	"github.com/kviz48/pocket-client-bot/internal/storage"
	"github.com/kviz48/pocket-client-bot/internal/storage/boltdb"
	"github.com/kviz48/pocket-client-bot/internal/telegram"
	"github.com/zhashkevych/go-pocket-sdk"
	"log"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	tgClient, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		log.Fatalf("telegram client init failed: %s", err)
	}

	tgClient.Debug = cfg.Debug

	pocketClient, err := pocket.NewClient(cfg.PocketToken)
	if err != nil {
		log.Fatalf("pocket client init failed: %s", err)
	}

	db, err := initDB(cfg)

	if err != nil {
		log.Fatalf("database init: %s", err)
	}
	tokenStorage := boltdb.NewTokenStorage(db)

	authServer := server.NewAuthServer(pocketClient, tokenStorage, cfg.BotURL)

	go func() {
		err := authServer.Start()
		if err != nil {
			log.Fatalf("authorization server failed: %s", err)
		}
	}()

	tgBot, err := telegram.NewBot(tgClient, pocketClient, cfg.AuthServerURL, tokenStorage, cfg.Messages)
	if err != nil {
		log.Fatalf("bot init failed: %s", err)
	}

	if err := tgBot.Start(); err != nil {
		log.Fatalf("cannot start bot: %s", err)
	}

	//	TODO: graceful shutdown
}

func initDB(cfg *config.Config) (*bolt.DB, error) {
	db, err := bolt.Open(cfg.DBFile, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(storage.RequestTokens))
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte(storage.AccessTokens))
		if err != nil {
			return err
		}

		return nil
	})

	return db, err
}
