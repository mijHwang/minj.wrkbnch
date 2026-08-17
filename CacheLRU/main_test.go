package main

import (
	"testing"
	"time"
)

func TestLRUCacheEviction(t *testing.T) {
	cache := NewLRUCache(2)

	cache.Set("a", "1", time.Minute)
	cache.Set("b", "2", time.Minute)
	cache.Set("c", "3", time.Minute) // capacity is 2 — this should evict "a"

	if _, ok := cache.Get("a"); ok {
		t.Errorf("expected 'a' to be evicted, but it was still found")
	}

	if val, ok := cache.Get("b"); !ok || val != "2" {
		t.Errorf("expected 'b' present with value '2', got value=%q ok=%v", val, ok)
	}

	if val, ok := cache.Get("c"); !ok || val != "3" {
		t.Errorf("expected 'c' present with value '3', got value=%q ok=%v", val, ok)
	}

}

func TestLRUCacheExpiry(t *testing.T) {
	cache := NewLRUCache(2)

	cache.Set("x", "temp", 50*time.Millisecond)

	if val, ok := cache.Get("x"); !ok || val != "temp" {
		t.Errorf("expected 'x' present immediately after Set, got value=%q ok=%v", val, ok)
	}

	time.Sleep(100 * time.Millisecond) // wait past the 50ms TTL

	if _, ok := cache.Get("x"); ok {
		t.Errorf("expected 'x' to be expired, but it was still found")
	}
}
