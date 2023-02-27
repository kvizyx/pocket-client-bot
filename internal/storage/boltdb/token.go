package boltdb

import (
	"errors"
	"github.com/boltdb/bolt"
	"github.com/kviz48/pocket-client-bot/internal/storage"
	"strconv"
)

type TokenStorage struct {
	db *bolt.DB
}

func NewTokenStorage(db *bolt.DB) *TokenStorage {
	return &TokenStorage{db: db}
}

func (s *TokenStorage) Save(chatID int64, token string, bucket storage.Bucket) error {
	err := s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		err := b.Put(intToBytes(chatID), []byte(token))
		return err
	})

	return err
}

func (s *TokenStorage) Get(chatID int64, bucket storage.Bucket) (string, error) {
	var token string
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		data := b.Get(intToBytes(chatID))
		token = string(data)
		return nil
	})

	if err != nil {
		return "", err
	}

	if token == "" {
		return "", errors.New("token not found")
	}

	return string(token), err
}

func intToBytes(v int64) []byte {
	return []byte(strconv.FormatInt(v, 10))
}
