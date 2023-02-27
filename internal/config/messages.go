package config

import (
	"errors"
	"github.com/spf13/viper"
)

type MessagesConfig struct {
	Errors    Errors
	Responses Responses
}

type Errors struct {
	Default      string `mapstructure:"default"`
	InvalidLink  string `mapstructure:"invalid_link"`
	Unauthorized string `mapstructure:"unauthorized"`
}

type Responses struct {
	Start             string `mapstructure:"start"`
	UnknownCommand    string `mapstructure:"unknown_command"`
	LinkSaveSuccess   string `mapstructure:"link_saved"`
	AlreadyAuthorized string `mapstructure:"already_authorized"`
}

func NewMessagesConfig() (*MessagesConfig, error) {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config.messages")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	var cfg MessagesConfig

	err := viper.UnmarshalKey("errors", &cfg.Errors)
	if err != nil {
		return nil, err
	}

	err = viper.UnmarshalKey("responses", &cfg.Responses)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
