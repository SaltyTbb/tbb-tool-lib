package datastructures

import (
	"testing"
)

func TestLRUCache_Remove(t *testing.T) {
	t.Run("remove existing key", func(t *testing.T) {
		// Setup
		cache := NewLruCache[string, int](5)
		cache.Put("key1", 1)
		cache.Put("key2", 2)
		cache.Put("key3", 3)

		// Initial state check
		if cache.Len() != 3 {
			t.Errorf("Expected cache length of 3, got %d", cache.Len())
		}

		// Execute
		cache.Remove("key2")

		// Verify removal
		if cache.Len() != 2 {
			t.Errorf("Expected cache length of 2 after removal, got %d", cache.Len())
		}

		// Verify key is actually gone
		_, ok := cache.Get("key2")
		if ok {
			t.Error("Key 'key2' should not exist after removal")
		}

		// Verify other keys still exist
		val1, ok1 := cache.Get("key1")
		if !ok1 || val1 != 1 {
			t.Errorf("Key 'key1' should still exist with value 1, got %v, exists: %v", val1, ok1)
		}

		val3, ok3 := cache.Get("key3")
		if !ok3 || val3 != 3 {
			t.Errorf("Key 'key3' should still exist with value 3, got %v, exists: %v", val3, ok3)
		}
	})

	t.Run("remove non-existent key", func(t *testing.T) {
		// Setup
		cache := NewLruCache[string, int](5)
		cache.Put("key1", 1)
		cache.Put("key2", 2)

		// Execute - this should not panic
		cache.Remove("key3")

		// Verify state is unchanged
		if cache.Len() != 2 {
			t.Errorf("Expected cache length to remain 2, got %d", cache.Len())
		}

		// Verify existing keys are intact
		val1, ok1 := cache.Get("key1")
		if !ok1 || val1 != 1 {
			t.Errorf("Key 'key1' should still exist with value 1, got %v, exists: %v", val1, ok1)
		}
	})

	t.Run("remove and add back", func(t *testing.T) {
		// Setup
		cache := NewLruCache[string, int](5)
		cache.Put("key1", 1)
		cache.Put("key2", 2)

		// Remove
		cache.Remove("key1")

		// Verify removal
		_, ok := cache.Get("key1")
		if ok {
			t.Error("Key 'key1' should be removed")
		}

		// Add back with different value
		cache.Put("key1", 100)

		// Verify addition
		val, ok := cache.Get("key1")
		if !ok {
			t.Error("Key 'key1' should exist after re-adding")
		}
		if val != 100 {
			t.Errorf("Expected value 100 for key 'key1', got %d", val)
		}
	})

	t.Run("capacity behavior after removal", func(t *testing.T) {
		// Setup a cache with capacity 3
		cache := NewLruCache[string, int](3)
		cache.Put("key1", 1)
		cache.Put("key2", 2)
		cache.Put("key3", 3)

		// Remove one item
		cache.Remove("key2")

		// Add two more - this will exceed capacity (which is still 3)
		// and will cause eviction of the least recently used key
		cache.Put("key4", 4)
		cache.Put("key5", 5)

		// key1 should be evicted as it's the least recently used
		_, ok1 := cache.Get("key1")
		if ok1 {
			t.Error("Key 'key1' should have been evicted as LRU")
		}

		// key2 was explicitly removed
		_, ok2 := cache.Get("key2")
		if ok2 {
			t.Error("Key 'key2' should have been removed")
		}

		// key3, key4, key5 should exist
		_, ok3 := cache.Get("key3")
		if !ok3 {
			t.Error("Key 'key3' should still exist")
		}

		_, ok4 := cache.Get("key4")
		if !ok4 {
			t.Error("Key 'key4' should exist")
		}

		_, ok5 := cache.Get("key5")
		if !ok5 {
			t.Error("Key 'key5' should exist")
		}

		// Verify length is at capacity
		if cache.Len() != 3 {
			t.Errorf("Expected cache length to be at capacity (3), got %d", cache.Len())
		}

		// Adding one more should evict the least recently used (key3, since we just accessed key4 and key5)
		cache.Put("key6", 6)

		_, ok3 = cache.Get("key3")
		if ok3 {
			t.Error("Key 'key3' should have been evicted as LRU")
		}
	})

	t.Run("remove all items", func(t *testing.T) {
		// Setup
		cache := NewLruCache[string, int](5)
		cache.Put("key1", 1)
		cache.Put("key2", 2)
		cache.Put("key3", 3)

		// Remove all
		cache.Remove("key1")
		cache.Remove("key2")
		cache.Remove("key3")

		// Verify empty
		if cache.Len() != 0 {
			t.Errorf("Expected empty cache, got length %d", cache.Len())
		}

		// Add new items
		cache.Put("keyA", 10)
		cache.Put("keyB", 20)

		// Verify new items
		valA, okA := cache.Get("keyA")
		if !okA || valA != 10 {
			t.Errorf("Expected keyA=10, got %v, exists: %v", valA, okA)
		}

		valB, okB := cache.Get("keyB")
		if !okB || valB != 20 {
			t.Errorf("Expected keyB=20, got %v, exists: %v", valB, okB)
		}
	})
}
