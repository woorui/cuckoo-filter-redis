# cuckoo-filter-redis

Concurrent, persistable, stand-alone <https://github.com/linvon/cuckoo-filter> .

## example

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	filter "github.com/woorui/cuckoo-filter-redis"
)

func main() {
	kv := filter.NewRedisKV(redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"}), "YOUR_REDIS_KEY")

	filter, close, err := filter.NewFilter(context.TODO(), kv, 24*time.Hour, filter.MemNewFilter(4, 9, 3900, filter.TableTypePacked))
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

// Output:
// The filter contain A 
// The filter size is 1
```

Note: This cuckoo-filter with redis-backed can't working in A distributed environment.
