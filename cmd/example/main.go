package main

import (
	"fmt"

	"github.com/SaltyTbb/tbb-tool-lib/pkg/datastructures"
)

func main() {
	cache := datastructures.NewLruCache[string, string](10)
	cache.Put("key", "value")
	value, ok := cache.Get("key")
	if !ok {
		fmt.Println("key not found")
	}
	fmt.Println(value)
	fmt.Println(cache)
	cache.Remove("key")
	cache.Put("key2", "value2")
	_, ok = cache.Get("key2")
	if !ok {
		fmt.Println("key2 not found")
	}
	fmt.Println(cache.String())
	fmt.Println(cache)
}
