# minerva
[![Build Status](https://travis-ci.org/AdhityaRamadhanus/minerva.svg?branch=master)](https://travis-ci.org/AdhityaRamadhanus/minerva)

configuration management using redis

<p>
  <a href="#installation">Installation |</a>
  <a href="#usage">Usage |</a>
  <a href="#todo">Todo |</a>
  <a href="#licenses">License</a>
  <br><br>
  <blockquote>
	minerva is remote config library using redis as configuration management
  </blockquote>
</p>

Installation
----------- 
* go get github.com/AdhityaRamadhanus/minerva

Usage
----------------
* See examples for more details
```go
package main

import (
	"log"
	"os"
	"time"

	"github.com/AdhityaRamadhanus/minerva"
	redis "github.com/AdhityaRamadhanus/minerva/redis"
)

func main() {
    // Create your service provider, currently only redis is implemented
	redisService := redis.NewRedisService(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})

    minerva := minerva.New(redisService)
    // you can watch any change on the config by calling Watch function
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
```

Todo
----------------
* ~~Dynamic prefix key~~
* Bootstrap config value
* CLI to manage configuration
* ~~More robust and proper error handling~~
* Add debouncing in getting key value on parsing key event
* Add test
* Add CI

License
----

MIT © [Adhitya Ramadhanus]

