package limit

import (
	"fmt"
	"sync"

	"github.com/go-redis/redis"
	"github.com/mshindle/ratelimit"
)

func RunRequest() {
	client := redis.NewClient(&redis.Options{
		Addr:     "firefly.dev:6379",
		Password: "",
		DB:       0,
	})
	pong, err := client.Ping().Result()
	fmt.Println(pong, err)

	var wg sync.WaitGroup
	_ = &ratelimit.RequestRate{ReplenishRate: 2, Capacity: 2, Client: client}

	for x := 0; x < 4; x++ {
		//initRequest(x, &wg, rr)
	}
	wg.Wait()
}
