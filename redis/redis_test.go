package redis

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
)

var redisService *RedisService

func TestMain(m *testing.M) {
	// init test condition
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   1,
	})

	if err := redisClient.ConfigSet("notify-keyspace-events", "KEA").Err(); err != nil {
		fmt.Printf("could not set config, %v", err)
		os.Exit(1)
	}

	if err := redisClient.FlushAll().Err(); err != nil {
		fmt.Printf("could not flush all keys, %v", err)
		os.Exit(1)
	}

	redisState := []struct {
		Key   string
		Value string
	}{
		{
			Key:   "test-1",
			Value: "true",
		},
		{
			Key:   "test-2",
			Value: "true",
		},
		{
			Key:   "test-3",
			Value: "true",
		},
	}

	for _, state := range redisState {
		if err := redisClient.Set(state.Key, state.Value, 0).Err(); err != nil {
			fmt.Printf("could not set key value, %v", err)
			os.Exit(1)
		}
	}

	redisService = New(&Options{
		Addr: "localhost:6379",
		DB:   1,
	})

	os.Exit(m.Run())
}

func TestGetKey(t *testing.T) {
	testCases := []struct {
		Key           string
		ExpectedValue string
	}{
		{
			Key:           "test-1",
			ExpectedValue: "true",
		},
		{
			Key:           "test-2",
			ExpectedValue: "true",
		},
		{
			Key:           "test-3",
			ExpectedValue: "true",
		},
		{
			Key:           "_oskdoaskdoaksdokasd",
			ExpectedValue: "",
		},
	}

	for _, test := range testCases {
		actualValue := redisService.Get(test.Key)
		assert.Equal(t, test.ExpectedValue, actualValue)
	}
}

func TestSetKey(t *testing.T) {
	testCases := []struct {
		Key   string
		Value string
	}{
		{
			Key:   "test-1",
			Value: "false",
		},
		{
			Key:   "test-2",
			Value: "false",
		},
	}

	for _, test := range testCases {
		assert.NoError(t, redisService.Set(test.Key, test.Value))
	}
}

func TestWatchKey(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	prefixKey := "config"
	keyEventChannel, err := redisService.Watch(ctx, prefixKey)
	assert.NoError(t, err)

	testCases := []struct {
		Key   string
		Value string
	}{
		{
			Key:   "test-1",
			Value: "true",
		},
		{
			Key:   "test-2",
			Value: "true",
		},
	}

	eventCounter := 0

	go func(ctx context.Context) {
		for {
			select {
			case _ = <-keyEventChannel:
				eventCounter++
			case <-ctx.Done():
				return
			}
		}
	}(ctx)

	for _, test := range testCases {
		redisKey := fmt.Sprintf("%s:%s", prefixKey, test.Key)
		redisService.Set(redisKey, test.Value)
	}

	// give time for key event to reach watcher
	time.Sleep(1 * time.Second)
	cancel()

	assert.Equal(t, 2, eventCounter)
}
