package storage

type Bucket string

const (
	RequestTokens Bucket = "request_tokens"
	AccessTokens  Bucket = "access_tokens"
)

type TokenStorage interface {
	Save(chatID int64, token string, bucket Bucket) error
	Get(chatID int64, bucket Bucket) (string, error)
}
