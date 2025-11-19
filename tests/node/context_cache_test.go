package node_test

import (
	"testing"

	"github.com/mautops/approval-kit/internal/node"
)

// TestContextCacheGetSet 测试上下文缓存的基本操作
func TestContextCacheGetSet(t *testing.T) {
	cache := node.NewContextCache()

	// 设置缓存数据
	cache.Set("key1", "value1")
	cache.Set("key2", 123)

	// 获取缓存数据
	value1, exists := cache.Get("key1")
	if !exists {
		t.Error("Cache key1 not found")
	}
	if value1 != "value1" {
		t.Errorf("Cache value1 = %v, want %v", value1, "value1")
	}

	value2, exists := cache.Get("key2")
	if !exists {
		t.Error("Cache key2 not found")
	}
	if value2 != 123 {
		t.Errorf("Cache value2 = %v, want %v", value2, 123)
	}
}

// TestContextCacheNotFound 测试缓存中不存在的键
func TestContextCacheNotFound(t *testing.T) {
	cache := node.NewContextCache()

	// 尝试获取不存在的键
	value, exists := cache.Get("non-existent")
	if exists {
		t.Error("Cache should not contain non-existent key")
	}
	if value != nil {
		t.Errorf("Cache value for non-existent key = %v, want nil", value)
	}
}

// TestContextCacheConcurrent 测试并发访问缓存
func TestContextCacheConcurrent(t *testing.T) {
	cache := node.NewContextCache()

	// 并发设置和获取
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(idx int) {
			key := "key" + string(rune('0'+idx))
			cache.Set(key, idx)
			value, exists := cache.Get(key)
			if !exists {
				t.Errorf("Cache key %s not found", key)
			}
			if value != idx {
				t.Errorf("Cache value for %s = %v, want %v", key, value, idx)
			}
			done <- true
		}(i)
	}

	// 等待所有 goroutine 完成
	for i := 0; i < 10; i++ {
		<-done
	}
}

// TestContextCacheNil 测试 nil 缓存
func TestContextCacheNil(t *testing.T) {
	var cache *node.ContextCache = nil

	// 对 nil 缓存的操作应该不会 panic
	value, exists := cache.Get("key")
	if exists {
		t.Error("Nil cache should not return exists=true")
	}
	if value != nil {
		t.Errorf("Nil cache value = %v, want nil", value)
	}

	// Set 操作应该不会 panic
	cache.Set("key", "value")
}

// TestContextCacheOverwrite 测试覆盖缓存值
func TestContextCacheOverwrite(t *testing.T) {
	cache := node.NewContextCache()

	// 设置初始值
	cache.Set("key", "value1")

	// 覆盖值
	cache.Set("key", "value2")

	// 验证新值
	value, exists := cache.Get("key")
	if !exists {
		t.Error("Cache key not found after overwrite")
	}
	if value != "value2" {
		t.Errorf("Cache value = %v, want %v", value, "value2")
	}
}

