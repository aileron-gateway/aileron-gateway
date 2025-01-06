package kvs

import (
	"context"
	"sync"
	"time"
)

// staticCtx is the context that won't be modified.
var staticCtx = context.TODO()

// MapKVS is the key-value store that uses built-in map
// as the data store.
// This is not for production use but only for development.
// MapKVS implements KVS interface.
type MapKVS[K comparable, V any] struct {
	Store map[K]V
	timer map[K]*time.Timer
	mu    sync.RWMutex
}

func (m *MapKVS[K, V]) Open(_ context.Context) error {
	if m.Store == nil {
		m.Store = map[K]V{}
		m.timer = map[K]*time.Timer{}
	}
	return nil
}

func (m *MapKVS[K, V]) Close(_ context.Context) error {
	return nil
}

func (m *MapKVS[K, V]) Get(_ context.Context, key K) (V, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if v, ok := m.Store[key]; ok {
		return v, nil
	}
	return *new(V), Nil
}

func (m *MapKVS[K, V]) Set(_ context.Context, key K, value V) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Store[key] = value
	return nil
}

func (m *MapKVS[K, V]) Delete(_ context.Context, key K) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.Store, key)
	if t, ok := m.timer[key]; ok {
		t.Stop()
	}
	delete(m.timer, key)
	return nil
}

func (m *MapKVS[K, V]) Exists(_ context.Context, key K) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, ok := m.Store[key]
	return ok
}

func (m *MapKVS[K, V]) SetWithTTL(ctx context.Context, key K, value V, ttl time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Store[key] = value

	if ttl > 0 {
		timer, ok := m.timer[key]
		if ok {
			timer.Reset(ttl)
			return nil
		}
		m.timer[key] = time.AfterFunc(ttl, func() { //nolint:contextcheck // Function `SetWithTTL$1` should pass the context parameter
			_ = m.Delete(staticCtx, key) // No need to check the returned error.
		})
	}

	return nil
}
