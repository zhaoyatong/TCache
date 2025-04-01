package cache

import "container/list"

// Value 存储记录的Value是一个接口类型
type Value interface {
	Len() int // 返回所占内存的大小
}

// LRUCache Cache LRU缓存结构体
type LRUCache struct {
	maxBytes  int64                         // 允许使用的最大内存
	usedBytes int64                         // 当前已使用的内存
	ll        *list.List                    // Go语言标准库实现的双向链表
	cache     map[string]*list.Element      // 键是字符串，值是双向链表中对应节点的指针
	OnEvicted func(key string, value Value) // 当记录删除时的回调函数，可以为nil
}

// NewLRUCache 创建一个LRUCache实例
func NewLRUCache(maxBytes int64, onEvicted func(string, Value)) *LRUCache {
	return &LRUCache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// 存储的记录结构
type entry struct {
	key   string
	value Value
}

// Get 获取记录
func (c *LRUCache) Get(key string) (value Value, ok bool) {
	if element, ok := c.cache[key]; ok {
		c.ll.MoveToFront(element) // 这里移到了队首
		kv := element.Value.(*entry)
		return kv.value, true
	}
	return
}

func (c *LRUCache) removeOldest() {
	element := c.ll.Back() // 取到队尾节点，从链表中删除。
	if element != nil {
		c.ll.Remove(element)
		kv := element.Value.(*entry)
		delete(c.cache, kv.key)
		c.usedBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

func (c *LRUCache) Set(key string, value Value) {
	// 如果键存在，则更新对应节点的值，并将该节点移到队首
	if element, ok := c.cache[key]; ok {
		c.ll.MoveToFront(element)
		kv := element.Value.(*entry)
		c.usedBytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		// 队首添加新节点 &entry{key, value}, 并字典中添加 key 和节点的映射关系
		element := c.ll.PushFront(&entry{key, value})
		c.cache[key] = element
		c.usedBytes += int64(len(key)) + int64(value.Len())
	}
	// 如果超过了设定的最大值 c.maxBytes，则移除最少访问的节点。
	for c.maxBytes != 0 && c.maxBytes < c.usedBytes {
		c.removeOldest()
	}
}
