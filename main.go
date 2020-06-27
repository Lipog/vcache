package main

import (
	"Vcache/vcache"
	"fmt"
	"log"
	"net/http"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func main() {
	vcache.NewGroup("scores", 2<<10, vcache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	addr := ":8080"
	peers := vcache.NewHTTPPool(addr)
	log.Println("vcache is running at", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}
