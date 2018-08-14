package redis

import (
	"context"
	"log"

	"github.com/AdhityaRamadhanus/minerva"
	"github.com/go-redis/redis"
)

type Options struct {
	Addr     string
	Password string
	DB       int
}

type RedisService struct {
	client *redis.Client
}

func NewRedisService(option *Options) *RedisService {
	client := redis.NewClient(&redis.Options{
		Addr:     option.Addr,
		Password: option.Password,
		DB:       option.DB,
	})

	return &RedisService{
		client: client,
	}
}

func NewRedisServiceWithClient(redisClient *redis.Client) *RedisService {
	return &RedisService{
		client: redisClient,
	}
}

func (r *RedisService) Get(key string) string {
	val, err := r.client.Get(key).Result()
	if err != nil {
		return ""
	}
	return val
}

func (r *RedisService) Set(key, value string) error {
	return r.client.Set(key, value, 0).Err()
}

func (r *RedisService) Watch(ctx context.Context) (chan minerva.KeyEvent, error) {
	log.Println("Watching")
	pubsub := r.client.PSubscribe("__key*__:config*")

	// Wait for confirmation that subscription is created before publishing anything.
	_, err := pubsub.Receive()
	if err != nil {
		return nil, err
	}

	// Go channel which receives messages.
	redisEventChannel := pubsub.Channel()
	keyEventChannel := make(chan minerva.KeyEvent, 10)
	go func() {
		defer pubsub.Close()
		for {
			select {
			case msg, ok := <-redisEventChannel:
				if !ok {
					break
				}
				keyEvent := parseKeyEvent(msg.Channel, msg.Payload)
				keyEvent.Value = r.Get("config:" + keyEvent.AffectedKey)
				keyEventChannel <- keyEvent
				log.Println("Redis Event:", msg.Channel, msg.Payload)
			case <-ctx.Done():
				return
			}
		}
	}()
	return keyEventChannel, nil
}
