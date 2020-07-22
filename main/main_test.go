package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func BenchmarkSplit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Splict()
	}
}
func Splict() {
	resp, _ := http.Get("http://localhost:9999/api?key=Tom")
	resp.Body.Close()
	bytes, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(bytes))
}