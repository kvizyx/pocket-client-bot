package server

import (
	"github.com/kviz48/pocket-client-bot/internal/storage"
	"github.com/kviz48/pocket-client-bot/internal/storage/boltdb"
	"github.com/zhashkevych/go-pocket-sdk"
	"net/http"
	"strconv"
)

type AuthServer struct {
	server      *http.Server
	pocket      *pocket.Client
	storage     *boltdb.TokenStorage
	redirectURL string
}

func NewAuthServer(pocket *pocket.Client, storage *boltdb.TokenStorage, redirectURL string) *AuthServer {
	return &AuthServer{
		pocket:      pocket,
		storage:     storage,
		redirectURL: redirectURL,
	}
}

func (s *AuthServer) Start() error {
	s.server = &http.Server{
		Addr:    ":80",
		Handler: s,
	}

	return s.server.ListenAndServe()
}

func (s *AuthServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	chatIDParam := r.URL.Query().Get("chat_id")
	if chatIDParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	chatID, err := strconv.ParseInt(chatIDParam, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	requestToken, err := s.storage.Get(chatID, storage.RequestTokens)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	accessToken, err := s.pocket.Authorize(r.Context(), requestToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = s.storage.Save(chatID, accessToken.AccessToken, storage.AccessTokens)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", s.redirectURL)
	w.WriteHeader(http.StatusMovedPermanently)
}
