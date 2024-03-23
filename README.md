
# Fooche

A naive caching library created to be used on [dcache](https://github.com/joaovictorsl/dcache).  



## Install

```bash
  go get github.com/joaovictorsl/fooche
```
    
## Usage

```go
package example

import (
	"fmt"
	"time"

	"github.com/joaovictorsl/fooche"
	"github.com/joaovictorsl/fooche/evict"
)

func main() {
	// Time in which expired data will be cleaned from cache
	cleanInterval := 10 * time.Second
	// Maps a data size (bytes) to its capacity
	// If 10 is mapped to 4 it means we'll be able to store 4 keys
    // with data size of 10 bytes or less
	sizeToCapMap := map[int]int{
		10: 4,
	}
	// Function to create an eviction policy based on storage capacity
	createLruPolicy := func(capacity int) evict.EvictionPolicy[string] {
		return evict.NewLRU[string](capacity)
	}
	// Creating caches
	// Only bounded caches use eviction policies
	c1 := fooche.NewCleanIntervalBounded(cleanInterval, sizeToCapMap, createLruPolicy)
	// c2 := fooche.NewCleanInterval(cleanInterval)
	// c3 := fooche.NewSimpleBounded(sizeToCapMap, createLruPolicy)
	// c4 := fooche.NewSimple()

	c1.Delete("key")
	c1.Get("key")
	c1.Has("key")
	c1.Set("key", []byte("data"), time.Second*30)
	fmt.Println(c1.String())
}
```

