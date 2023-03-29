package gocache

import (
	"hash/fnv"
	"sync"
)

type Container interface {
	Get(key string) string
	Set(key, value string)
	Del(key string)
}

type internalContainer struct {
	buckets []*bucket
}

type bucket struct {
	mu   sync.Mutex
	data map[string]string
}

func InitContainer(concurrency int) Container {
	size := 1
	for size < concurrency {
		size = concurrency << 1
	}
	container := internalContainer{
		buckets: make([]*bucket, size),
	}
	for i, _ := range container.buckets {
		container.buckets[i] = &bucket{data: make(map[string]string)}
	}
	return &container
}

func (c *internalContainer) Get(key string) string {
	index := spread(hashString(key)) % int32(len(c.buckets))
	bucket := c.buckets[index]
	return bucket.Get(key)
}

func (c *internalContainer) Set(key, value string) {
	index := spread(hashString(key)) % int32(len(c.buckets))
	bucket := c.buckets[index]
	bucket.Set(key, value)
}

func (c *internalContainer) Del(key string) {
	index := spread(hashString(key)) % int32(len(c.buckets))
	bucket := c.buckets[index]
	bucket.Del(key)
}

func (b *bucket) Get(key string) string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.data[key]
}

func (b *bucket) Set(key, value string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.data[key] = value
}

func (b *bucket) Del(key string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.data, key)
}

func spread(h int32) int32 {
	return (h ^ (h >> 16)) & int32(0x7fffffff)
}

func hashString(s string) int32 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(s))
	return int32(h.Sum32())
}
