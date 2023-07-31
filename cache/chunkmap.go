package cache

import (
	"fmt"
	"sync"
)

type ChunkMap map[string]map[int]string

var chunkMapCache map[string]map[int]string
var locker sync.RWMutex

func Get(key string) (map[int]string, bool) {
	locker.Lock()
	fmt.Println(chunkMapCache)
	v, ok := chunkMapCache[key]
	locker.Unlock()
	return v, ok
}

func Set(key string, index int, value string) {
	locker.Lock()
	m := chunkMapCache[key]
	if m == nil {
		m = make(map[int]string)
	}
	m[index] = value
	chunkMapCache[key] = m
	locker.Unlock()
}

func Del(key string) {
	locker.Lock()
	delete(chunkMapCache, key)
	locker.Unlock()
}

func init() {
	chunkMapCache = make(map[string]map[int]string)
}
