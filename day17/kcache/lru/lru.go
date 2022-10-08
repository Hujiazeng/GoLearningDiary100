package lru

import "container/list"

// Cache is a LRU cache.
type Cache struct {
	maxBytes   int64                    // 最大内存
	nbytes     int64                    // 当前已使用内存
	ll         *list.List               // 链表
	cache      map[string]*list.Element // 存放节点指针
	OnCallBack func(key string, value Value)
}

// 链表中的数据格式
type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

func New(maxBytes int64, OnCallBack func(string, Value)) *Cache {
	return &Cache{
		maxBytes:   maxBytes,
		ll:         list.New(),
		cache:      make(map[string]*list.Element),
		OnCallBack: OnCallBack,
	}
}

// 获取节点同时将节点移动到队列最前
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		// 检测节点类型
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// 缓存淘汰 删除最近最少访问的节点(队尾)
func (c *Cache) RemoveOldPoint() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		// 删除字典中的元素
		delete(c.cache, kv.key)

		// key和value的总长度
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnCallBack != nil {
			c.OnCallBack(kv.key, kv.value)
		}
	}
}

// 新增
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		// 键存在 更新值
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}

	// 超过最大值删除
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldPoint()
	}
}
