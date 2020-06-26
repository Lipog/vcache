package Vcache

import (
	lru2 "Vcache/lru"
	"testing"
)

type String string
func (d String) Len() int {
	return len(d)
}

func TestGet(t *testing.T) {
	lru := lru2.New(int64(0),nil)
	lru.Add("key1", String("123436"))
	if v, ok := lru.Get("key1"); !ok || string(v.(String)) != "123436" {
		t.Fatal("cache hit key1=123436 failed")
	}
	if _, ok :=  lru.Get("key2"); ok {
		t.Fatal("cache miss key2 failed")
	}
}
