package lru

import "testing"

type str string

func (s str) Len() int {
	return len(s)
}

func TestGet(t *testing.T) {
	lru := New(int64(0), nil)
	lru.Add("key1", str("1234"))
	if v, ok := lru.Get("key1"); !ok || string(v.(str)) != "1234" {
		t.Fatal("cache failed")
	}
	lru.Add("key1", str("666"))
	if v, ok := lru.Get("key1"); !ok || string(v.(str)) != "666" {
		t.Fatal("err")
	}

}

func TestRemove(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "key3"
	v1, v2, v3 := "value1", "value2", "value3"
	cap := len(k1 + k2 + v1 + v2)
	lru := New(int64(cap), nil)
	lru.Add(k1, str(v1))
	lru.Add(k2, str(v2))
	lru.Add(k3, str(v3))
	if _, ok := lru.Get("key1"); ok {
		t.Fatal("err")
	}

}
