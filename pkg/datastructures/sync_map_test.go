package datastructures

import (
	"sync"
	"testing"
	"time"
)

func TestSyncMap_BasicOperations(t *testing.T) {
	t.Run("set and get", func(t *testing.T) {
		m := NewSyncMap[string, int]()
		m.Set("key1", 42)

		value, ok := m.Get("key1")
		if !ok || value != 42 {
			t.Errorf("Expected (42, true), got (%v, %v)", value, ok)
		}
	})

	t.Run("missing key", func(t *testing.T) {
		// Special case: our implementation should block on missing keys
		// not returning zero values, so we'll test this in a separate test
	})

	t.Run("delete", func(t *testing.T) {
		m := NewSyncMap[string, int]()
		m.Set("key1", 42)
		m.Delete("key1")

		// In our implementation, Get would block forever if the key doesn't exist,
		// so we can't easily test for non-existence.
		// We can verify length though
		if m.Len() != 0 {
			t.Errorf("Expected length 0, got %d", m.Len())
		}
	})
}

func TestSyncMap_ThreadSafety(t *testing.T) {
	t.Run("concurrent reads and writes", func(t *testing.T) {
		m := NewSyncMap[int, int]()
		const numOperations = 1000
		const numGoroutines = 10

		var wg sync.WaitGroup
		wg.Add(numGoroutines * 2) // for readers and writers

		// Populate with initial data
		for i := 0; i < numOperations; i++ {
			m.Set(i, i*10)
		}

		// Concurrent readers
		for i := 0; i < numGoroutines; i++ {
			go func() {
				defer wg.Done()
				for j := 0; j < numOperations; j++ {
					val, _ := m.Get(j)
					if val != j*10 {
						t.Errorf("Expected %d, got %d", j*10, val)
					}
				}
			}()
		}

		// Concurrent writers
		for i := 0; i < numGoroutines; i++ {
			go func() {
				defer wg.Done()
				for j := 0; j < numOperations; j++ {
					m.Set(j, j*10) // Setting to the same value
				}
			}()
		}

		wg.Wait()
	})
}

func TestSyncMap_BlockingGet(t *testing.T) {
	t.Run("get waits for set", func(t *testing.T) {
		m := NewSyncMap[string, int]()
		const key = "test-key"
		const value = 42
		var getResult int
		var getSucceeded bool

		// Setup a goroutine to perform Get, which should block
		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			defer wg.Done()
			getResult, getSucceeded = m.Get(key)
		}()

		// Wait a moment to ensure Get has started
		time.Sleep(100 * time.Millisecond)

		// Set the value, which should unblock Get
		m.Set(key, value)

		// Use a timeout to ensure test doesn't hang forever
		done := make(chan struct{})
		go func() {
			wg.Wait()
			close(done)
		}()

		select {
		case <-done:
			// Test passed - Get completed
			if !getSucceeded || getResult != value {
				t.Errorf("Expected Get to return (%d, true), got (%d, %v)", value, getResult, getSucceeded)
			}
		case <-time.After(2 * time.Second):
			t.Fatal("Test timed out - Get did not complete after Set")
		}
	})

	t.Run("multiple gets wait for set", func(t *testing.T) {
		m := NewSyncMap[string, int]()
		const key = "multi-get-key"
		const value = 99

		// Launch 5 goroutines all waiting on the same key
		const numWaiters = 5
		results := make(chan int, numWaiters)

		var wg sync.WaitGroup
		wg.Add(numWaiters)

		for i := 0; i < numWaiters; i++ {
			go func() {
				defer wg.Done()
				val, _ := m.Get(key)
				results <- val
			}()
		}

		// Give goroutines time to start
		time.Sleep(100 * time.Millisecond)

		// Set the value once, which should unblock all getters
		m.Set(key, value)

		// Use a timeout to ensure test doesn't hang
		done := make(chan struct{})
		go func() {
			wg.Wait()
			close(done)
		}()

		select {
		case <-done:
			// All getters completed, check results
			close(results)
			count := 0
			for result := range results {
				if result != value {
					t.Errorf("Expected %d, got %d", value, result)
				}
				count++
			}
			if count != numWaiters {
				t.Errorf("Expected %d results, got %d", numWaiters, count)
			}
		case <-time.After(2 * time.Second):
			t.Fatal("Test timed out - Get calls did not complete after Set")
		}
	})
}
