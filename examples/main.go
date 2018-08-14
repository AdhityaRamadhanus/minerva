package main

import (
	"log"
	"os"
	"time"

	"github.com/AdhityaRamadhanus/minerva"
	redis "github.com/AdhityaRamadhanus/minerva/redis"
)

func main() {
	redisService := redis.New(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})

	minerva := minerva.NewWithOptions(redisService, minerva.Options{
		// default to config
		PrefixKey: "minerva-config",
	})
	minerva.Watch()
	for {
		if minerva.Get("is-maintenance") == "true" {
			log.Println("Server is in maintenance")
			minerva.Close()
			os.Exit(0)
		}
		log.Println("Serving Request")
		time.Sleep(1 * time.Second)
	}
}
