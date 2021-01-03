package skipmap

import (
	"reflect"
	"sync"
	"sync/atomic"
	"testing"
)

func TestSyncMap(t *testing.T) {
	m := NewInt64()
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		i := i
		wg.Add(1)
		go func() {
			m.Store(int64(i), int(i+1000))
			wg.Done()
		}()
	}
	wg.Wait()
	wg.Add(1)
	go func() {
		m.Delete(600)
		wg.Done()
	}()
	wg.Wait()
	wg.Add(1)
	var count int64
	go func() {
		m.Range(func(key int64, val interface{}) bool {
			atomic.AddInt64(&count, 1)
			return true
		})
		wg.Done()
	}()
	wg.Wait()

	val, ok := m.Load(500)
	if !ok || reflect.TypeOf(val).Kind().String() != "int" || val.(int) != 1500 {
		t.Fatal("1")
	}

	_, ok = m.Load(600)
	if ok {
		t.Fatal("2")
	}

	if m.Len() != 999 || int(count) != m.Len() {
		t.Fatal("3")
	}
}

func TestSyncMap2(t *testing.T) {
}
