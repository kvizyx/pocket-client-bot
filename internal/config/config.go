package config

import (
	"errors"
	"github.com/spf13/viper"
)

type Config struct {
	BotToken      string
	PocketToken   string
	AuthServerURL string
	BotURL        string
	DBFile        string
	Debug         bool

	Messages *MessagesConfig
}

func NewConfig() (*Config, error) {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config.main")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	var cfg Config

	err := viper.UnmarshalKey("db.file", &cfg.DBFile)
	if err != nil {
		return nil, err
	}

	err = viper.UnmarshalKey("bot.url", &cfg.BotURL)
	if err != nil {
		return nil, err
	}

	err = viper.UnmarshalKey("bot.secret_token", &cfg.BotToken)
	if err != nil {
		return nil, err
	}

	err = viper.UnmarshalKey("server.auth_url", &cfg.AuthServerURL)
	if err != nil {
		return nil, err
	}

	err = viper.UnmarshalKey("api.consumer_key", &cfg.PocketToken)
	if err != nil {
		return nil, err
	}

	err = viper.UnmarshalKey("app.debug", &cfg.Debug)
	if err != nil {
		return nil, err
	}

	messagesCfg, err := NewMessagesConfig()
	if err != nil {
		return nil, err
	}

	cfg.Messages = messagesCfg

	return &cfg, nil
}
