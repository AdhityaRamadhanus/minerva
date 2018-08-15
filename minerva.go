package minerva

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
)

// Options hold various configuration for minerva
type Options struct {
	PrefixKey string
}

var defaultOptions = Options{
	PrefixKey: "config",
}

// Minerva is the main struct that glue all components for creating remote config
type Minerva struct {
	remoteClient RemoteClient
	remoteConfig map[string]string
	cancelFunc   context.CancelFunc
	context      context.Context
	options      Options
}

// New create Minerva with default Options
func New(remoteClient RemoteClient) *Minerva {
	return &Minerva{
		remoteClient: remoteClient,
		remoteConfig: map[string]string{},
		options:      defaultOptions,
	}
}

// NewWithOptions create Minerva with custom options
func NewWithOptions(remoteClient RemoteClient, options Options) *Minerva {
	return &Minerva{
		remoteClient: remoteClient,
		remoteConfig: map[string]string{},
		options:      options,
	}
}

// Get a key value in remote config
func (m *Minerva) Get(key string) string {
	_, isKeyPresent := m.remoteConfig[key]
	remoteConfigKey := fmt.Sprintf("%s:%s", m.options.PrefixKey, key)
	if !isKeyPresent {
		m.remoteConfig[key] = m.remoteClient.Get(remoteConfigKey)
	}
	return m.remoteConfig[key]
}

// Watch all changes happening on remote config and apply it on local config
func (m *Minerva) Watch() error {
	ctx, cancel := context.WithCancel(context.Background())
	m.context = ctx
	m.cancelFunc = cancel

	keyEventChannel, err := m.remoteClient.Watch(ctx, m.options.PrefixKey)

	go func(ctx context.Context) {
		for {
			select {
			case keyEvent := <-keyEventChannel:
				m.remoteConfig[keyEvent.AffectedKey] = keyEvent.Value
			case <-ctx.Done():
				return
			}
		}
	}(ctx)

	return errors.Wrap(err, "Error in watching key event")
}

// Close watcher
func (m *Minerva) Close() {
	m.cancelFunc()
}
