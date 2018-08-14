package minerva

import (
	"context"
)

type RemoteClient interface {
	Get(key string) string
	Set(key string, value string) error
	Watch(ctx context.Context) (chan KeyEvent, error)
}
