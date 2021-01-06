package skipmap

import (
	"reflect"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
)

func TestStringMap(t *testing.T) {
	m := NewString()

	// Correctness.
	m.Store("123", "123")
	v, ok := m.Load("123")
	if !ok || v != "123" || m.Len() != 1 {
		t.Fatal("invalid")
	}

	m.Store("123", "456")
	v, ok = m.Load("123")
	if !ok || v != "456" || m.Len() != 1 {
		t.Fatal("invalid")
	}

	m.Store("123", 456)
	v, ok = m.Load("123")
	if !ok || v != 456 || m.Len() != 1 {
		t.Fatal("invalid")
	}

	m.Delete("123")
	_, ok = m.Load("123")
	if ok || m.Len() != 0 {
		t.Fatal("invalid")
	}

	_, ok = m.LoadOrStore("123", 456)
	if ok || m.Len() != 1 {
		t.Fatal("invalid")
	}

	v, ok = m.Load("123")
	if !ok || v != 456 || m.Len() != 1 {
		t.Fatal("invalid")
	}

	v, ok = m.LoadAndDelete("123")
	if !ok || v != 456 || m.Len() != 0 {
		t.Fatal("invalid")
	}

	_, ok = m.LoadOrStore("123", 456)
	if ok || m.Len() != 1 {
		t.Fatal("invalid")
	}

	m.LoadOrStore("456", 123)
	if ok || m.Len() != 2 {
		t.Fatal("invalid")
	}

	m.Range(func(key string, value interface{}) bool {
		if key == "123" {
			m.Store("123", 123)
		} else if key == "456" {
			m.LoadAndDelete("456")
		}
		return true
	})

	v, ok = m.Load("123")
	if !ok || v != 123 || m.Len() != 1 {
		t.Fatal("invalid")
	}

	// Concurrent.
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		i := i
		wg.Add(1)
		go func() {
			n := strconv.Itoa(i)
			m.Store(n, int(i+1000))
			wg.Done()
		}()
	}
	wg.Wait()
	var count2 int64
	m.Range(func(key string, value interface{}) bool {
		atomic.AddInt64(&count2, 1)
		return true
	})
	m.Delete("600")
	var count int64
	m.Range(func(key string, value interface{}) bool {
		atomic.AddInt64(&count, 1)
		return true
	})

	val, ok := m.Load("500")
	if !ok || reflect.TypeOf(val).Kind().String() != "int" || val.(int) != 1500 {
		t.Fatal("fail")
	}

	_, ok = m.Load("600")
	if ok {
		t.Fatal("fail")
	}

	if m.Len() != 999 || int(count) != m.Len() {
		t.Fatal("fail", m.Len(), count, count2)
	}
}
