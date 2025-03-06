package pokecache

import (
	"bytes"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestAddGet(t *testing.T) {
	const interval = 5 * time.Second
	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "https://example.com",
			val: []byte("testdata"),
		},
		{
			key: "https://example.com/path",
			val: []byte("moretestdata"),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			cache := NewCache(interval)
			cache.Add(c.key, c.val)
			val, ok := cache.Get(c.key)
			if !ok {
				t.Errorf("expected to find key")
				return
			}
			if string(val) != string(c.val) {
				t.Errorf("expected to find value")
				return
			}
		})
	}
}

func TestReapLoop(t *testing.T) {
	const baseTime = 5 * time.Millisecond
	const waitTime = baseTime + 5*time.Millisecond
	cache := NewCache(baseTime)
	cache.Add("https://example.com", []byte("testdata"))

	_, ok := cache.Get("https://example.com")
	if !ok {
		t.Errorf("expected to find key")
		return
	}

	time.Sleep(waitTime)

	_, ok = cache.Get("https://example.com")
	if ok {
		t.Errorf("expected to not find key")
		return
	}
}

func TestConcurrentAccess(t *testing.T) {
	const interval = 5 * time.Second
	cache := NewCache(interval)

	// Number of concurrent operations
	const workers = 100
	var wg sync.WaitGroup
	wg.Add(workers * 2) // For both reading and writing

	// Add a baseline entry that all goroutines will attempt to access
	cache.Add("shared-key", []byte("original-value"))

	// Launch multiple goroutines to read and write simultaneously
	for i := 0; i < workers; i++ {
		// Writer goroutine
		go func(id int) {
			defer wg.Done()
			key := fmt.Sprintf("key-%d", id)
			cache.Add(key, []byte(fmt.Sprintf("value-%d", id)))
		}(i)

		// Reader goroutine
		go func(id int) {
			defer wg.Done()
			key := "shared-key"
			val, found := cache.Get(key)
			if !found {
				t.Errorf("Worker %d: Expected to find shared key", id)
			}
			if string(val) != "original-value" {
				t.Errorf("Worker %d: Expected original value, got %s", id, string(val))
			}
		}(i)
	}

	wg.Wait()

	// Verify all entries were properly added
	for i := 0; i < workers; i++ {
		key := fmt.Sprintf("key-%d", i)
		val, found := cache.Get(key)
		if !found {
			t.Errorf("Failed to find key: %s", key)
			continue
		}
		expectedVal := fmt.Sprintf("value-%d", i)
		if string(val) != expectedVal {
			t.Errorf("Expected %s, got %s for key %s", expectedVal, string(val), key)
		}
	}
}

func TestSpecialCases(t *testing.T) {
	const interval = 5 * time.Second

	testCases := []struct {
		name string
		key  string
		val  []byte
	}{
		{"Empty key", "", []byte("some value")},
		{"Empty value", "some-key", []byte("")},
		{"Nil value", "nil-key", nil},
		{"Unicode key", "こんにちは", []byte("hello")},
		{"Unicode value", "hello", []byte("こんにちは")},
		{"Very long key", strings.Repeat("x", 10000), []byte("value")},
		{"Very long value", "key", bytes.Repeat([]byte("x"), 10000)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cache := NewCache(interval)

			// This shouldn't panic
			cache.Add(tc.key, tc.val)

			val, found := cache.Get(tc.key)
			if !found {
				t.Fatal("Failed to find key")
			}

			if !bytes.Equal(val, tc.val) {
				t.Fatalf("Values don't match. Expected %v, got %v", tc.val, val)
			}
		})
	}
}
