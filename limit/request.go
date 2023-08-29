package limit

import (
	"context"
	"fmt"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/mshindle/ratelimit/rediskv"
)

func RunRequest(ctx context.Context) {
	client := redis.NewClient(&redis.Options{
		Addr:     "firefly.dev:6379",
		Password: "",
		DB:       0,
	})
	pong, err := client.Ping(ctx).Result()
	fmt.Println(pong, err)

	var wg sync.WaitGroup
	_ = rediskv.NewRateLimiter(client)
	//ratelimit.RequestRate{ReplenishRate: 2, Capacity: 2, Client: client}

	for x := 0; x < 4; x++ {
		//initRequest(x, &wg, rr)
	}
	wg.Wait()
}
