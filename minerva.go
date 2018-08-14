package minerva

import (
	"context"
	"log"
)

type Minerva struct {
	remoteClient RemoteClient
	remoteConfig map[string]string
	cancelFunc   context.CancelFunc
	context      context.Context
}

func New(remoteClient RemoteClient) *Minerva {
	return &Minerva{
		remoteClient: remoteClient,
		remoteConfig: map[string]string{},
	}
}

func (m *Minerva) Get(key string) string {
	_, isKeyPresent := m.remoteConfig[key]
	if !isKeyPresent {
		m.remoteConfig[key] = m.remoteClient.Get("config:" + key)
	}
	return m.remoteConfig[key]
}

// TODO: Add golang context
func (m *Minerva) Watch() error {
	ctx, cancel := context.WithCancel(context.Background())
	m.context = ctx
	m.cancelFunc = cancel

	keyEventChannel, err := m.remoteClient.Watch(ctx)

	go func(ctx context.Context) {
		for {
			select {
			case keyEvent := <-keyEventChannel:
				m.remoteConfig[keyEvent.AffectedKey] = keyEvent.Value
				log.Println("Minerva Event:", keyEvent.Type, keyEvent.AffectedKey, keyEvent.Value)
			case <-ctx.Done():
				return
			}
		}
	}(ctx)

	return err
}

func (m *Minerva) Close() {
	m.cancelFunc()
}
