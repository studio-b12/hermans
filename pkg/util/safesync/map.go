package safesync

import (
	"sync"
)

type Map[TKey any, TValue any] struct {
	inner sync.Map
}

func (t *Map[TKey, TValue]) Load(key TKey) (v TValue, ok bool) {
	vAny, ok := t.inner.Load(key)
	if !ok {
		return v, false
	}
	return vAny.(TValue), true
}

func (t *Map[TKey, TValue]) Store(key TKey, value TValue) {
	t.inner.Store(key, value)
}

func (t *Map[TKey, TValue]) Clear() {
	t.inner.Clear()
}
