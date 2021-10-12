package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/linvon/cuckoo-filter"
	filter "github.com/woorui/cuckoo-filter-redis"
)

func main() {
	kv := filter.NewRedisKV(redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"}), "YOUR_REDIS_KEY")

	filter, close, err := filter.NewFilter(context.TODO(), kv, 2*time.Second, cuckoo.NewFilter(4, 9, 3900, cuckoo.TableTypePacked))
	if err != nil {
		panic(err)
	}
	defer close()

	a := []byte("A")
	filter.Add(a)

	ok, _ := filter.Contain(a)
	if ok {
		fmt.Printf("The filter contain %s \n", string(a))
	}

	fmt.Printf("The filter size is %d \n", filter.Size())
}
