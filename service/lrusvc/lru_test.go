package lrusvc

import (
	"fmt"
	"testing"

	lru "github.com/hashicorp/golang-lru"
)

func TestLru(t *testing.T) {
	// 创建新的LRU缓存，容量为2
	cache, err := lru.New(2)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 添加键值对到缓存
	cache.Add("key1", "value1")
	cache.Add("key2", "value2")
	cache.Add("key3", "value3") // LRU缓存容量为2，新加入的元素将替换最近最少使用的元素

	// 获取缓存中的值
	value, ok := cache.Get("key1")
	if ok {
		fmt.Println("Value:", value)
	} else {
		fmt.Println("Key not found")
	}

	value, ok = cache.Get("key2")
	if ok {
		fmt.Println("Value:", value)
	} else {
		fmt.Println("Key not found")
	}

	value, ok = cache.Get("key3")
	if ok {
		fmt.Println("Value:", value)
	} else {
		fmt.Println("Key not found")
	}
}
