package minerva

import (
	"context"
)

// RemoteClient is interface to interact with remote config
type RemoteClient interface {
	Get(key string) string
	Set(key string, value string) error
	Watch(ctx context.Context, prefixKey string) (chan KeyEvent, error)
}
