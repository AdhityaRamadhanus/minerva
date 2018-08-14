package redis

import (
	"context"
	"fmt"

	"github.com/AdhityaRamadhanus/minerva"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

type Options struct {
	Addr     string
	Password string
	DB       int
}

type RedisService struct {
	client *redis.Client
}

func New(option *Options) *RedisService {
	client := redis.NewClient(&redis.Options{
		Addr:     option.Addr,
		Password: option.Password,
		DB:       option.DB,
	})

	return &RedisService{
		client: client,
	}
}

func NewWithClient(redisClient *redis.Client) *RedisService {
	return &RedisService{
		client: redisClient,
	}
}

func (r *RedisService) Get(key string) string {
	val, err := r.client.Get(key).Result()
	// supress error and return empty string instead
	if err != nil {
		return ""
	}
	return val
}

func (r *RedisService) Set(key, value string) error {
	return r.client.Set(key, value, 0).Err()
}

func (r *RedisService) Watch(ctx context.Context, prefixKey string) (chan minerva.KeyEvent, error) {
	keyspace := fmt.Sprintf("__key*__:%s*", prefixKey)
	pubsub := r.client.PSubscribe(keyspace)

	// Wait for confirmation that subscription is created before publishing anything.
	_, err := pubsub.Receive()
	if err != nil {
		errMessage := fmt.Sprintf("error in creating subcription to keyspace %s", keyspace)
		return nil, errors.Wrap(err, errMessage)
	}

	// redis key event channel
	redisEventChannel := pubsub.Channel()
	// exposed to Minerva, only contains KeyEvent
	keyEventChannel := make(chan minerva.KeyEvent, 10)
	go func() {
		defer pubsub.Close()
		for {
			select {
			case msg, ok := <-redisEventChannel:
				if !ok {
					break
				}
				keyEvent := minerva.KeyEvent{
					AffectedKey: parseKeyEvent(prefixKey, msg.Channel),
					Type:        msg.Payload,
				}
				redisKey := fmt.Sprintf("%s:%s", prefixKey, keyEvent.AffectedKey)
				// due to redis only provide the event and the affected key
				// we need to fetch the value
				keyEvent.Value = r.client.Get(redisKey).Val()
				keyEventChannel <- keyEvent
			case <-ctx.Done():
				return
			}
		}
	}()

	return keyEventChannel, nil
}
