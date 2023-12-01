package skipmap

import (
	"math"
	"math/rand"
	"reflect"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/zhangyunhao116/fastrand"
)

func TestTyped(t *testing.T) {
	testSkipMapInt(t, func() anyskipmap[int] { return NewInt[any]() })
	testSkipMapIntDesc(t, func() anyskipmap[int] { return NewIntDesc[any]() })
	testSkipMapString(t, func() anyskipmap[string] { return NewString[any]() })
	testSkipMapString(t, func() anyskipmap[string] { return NewStringDesc[any]() })
	testSkipMapString(t, func() anyskipmap[string] { return NewStringFast[any]() })
	testSyncMapSuiteInt64(t, func() anyskipmap[int64] { return NewInt64[any]() })
}

func TestOrdered(t *testing.T) {
	testSkipMapInt(t, func() anyskipmap[int] { return New[int, any]() })
	testSkipMapIntDesc(t, func() anyskipmap[int] { return NewDesc[int, any]() })
	testSkipMapString(t, func() anyskipmap[string] { return New[string, any]() })
	testSyncMapSuiteInt64(t, func() anyskipmap[int64] { return New[int64, any]() })
}

func TestFunc(t *testing.T) {
	testSkipMapInt(t, func() anyskipmap[int] { return NewFunc[int, any](func(a, b int) bool { return a < b }) })
}

type anyskipmap[T any] interface {
	Store(key T, value any)
	Load(key T) (any, bool)
	Delete(key T) bool
	LoadAndDelete(key T) (any, bool)
	LoadOrStore(key T, value any) (any, bool)
	LoadOrStoreLazy(key T, f func() any) (any, bool)
	Range(f func(key T, value any) bool)
	Len() int
}

func testSkipMapInt(t *testing.T, newset func() anyskipmap[int]) {
	m := newset()

	// Correctness.
	m.Store(123, "123")
	v, ok := m.Load(123)
	if !ok || v != "123" || m.Len() != 1 {
		t.Fatal("invalid")
	}

	m.Store(123, "456")
	v, ok = m.Load(123)
	if !ok || v != "456" || m.Len() != 1 {
		t.Fatal("invalid")
	}

	m.Store(123, 456)
	v, ok = m.Load(123)
	if !ok || v != 456 || m.Len() != 1 {
		t.Fatal("invalid")
	}

	m.Delete(123)
	v, ok = m.Load(123)
	if ok || m.Len() != 0 || v != nil {
		t.Fatal("invalid")
	}

	v, loaded := m.LoadOrStore(123, 456)
	if loaded || v != 456 || m.Len() != 1 {
		t.Fatal("invalid")
	}

	v, loaded = m.LoadOrStore(123, 789)
	if !loaded || v != 456 || m.Len() != 1 {
		t.Fatal("invalid")
	}

	v, ok = m.Load(123)
	if !ok || v != 456 || m.Len() != 1 {
		t.Fatal("invalid")
	}

	v, ok = m.LoadAndDelete(123)
	if !ok || v != 456 || m.Len() != 0 {
		t.Fatal("invalid")
	}

	_, ok = m.LoadOrStore(123, 456)
	if ok || m.Len() != 1 {
		t.Fatal("invalid")
	}

	m.LoadOrStore(456, 123)
	if ok || m.Len() != 2 {
		t.Fatal("invalid")
	}

	m.Range(func(key int, _ interface{}) bool {
		if key == 123 {
			m.Store(123, 123)
		} else if key == 456 {
			m.LoadAndDelete(456)
		}
		return true
	})

	v, ok = m.Load(123)
	if !ok || v != 123 || m.Len() != 1 {
		t.Fatal("invalid")
	}

	// Concurrent.
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		i := i
		wg.Add(1)
		go func() {
			m.Store(i, int(i+1000))
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
		m.Range(func(_ int, _ interface{}) bool {
			atomic.AddInt64(&count, 1)
			return true
		})
		wg.Done()
	}()
	wg.Wait()

	val, ok := m.Load(500)
	if !ok || reflect.TypeOf(val).Kind().String() != "int" || val.(int) != 1500 {
		t.Fatal("fail")
	}

	_, ok = m.Load(600)
	if ok {
		t.Fatal("fail")
	}

	if m.Len() != 999 || int(count) != m.Len() {
		t.Fatal("fail")
	}
	// Correctness 2.
	var m1 sync.Map
	m2 := newset()
	var v1, v2 interface{}
	var ok1, ok2 bool
	for i := 0; i < 100000; i++ {
		rd := int(fastrand.Uint32n(10))
		r1, r2 := int(fastrand.Uint32n(100)), int(fastrand.Uint32n(100))
		if rd == 0 {
			m1.Store(r1, r2)
			m2.Store(r1, r2)
		} else if rd == 1 {
			v1, ok1 = m1.LoadAndDelete(r1)
			v2, ok2 = m2.LoadAndDelete(r1)
			if ok1 != ok2 || v1 != v2 {
				t.Fatal(rd, v1, ok1, v2, ok2)
			}
		} else if rd == 2 {
			v1, ok1 = m1.LoadOrStore(r1, r2)
			v2, ok2 = m2.LoadOrStore(r1, r2)
			if ok1 != ok2 || v1 != v2 {
				t.Fatal(rd, v1, ok1, v2, ok2, "input -> ", r1, r2)
			}
		} else if rd == 3 {
			m1.Delete(r1)
			m2.Delete(r1)
		} else if rd == 4 {
			m2.Range(func(key int, value interface{}) bool {
				v, ok := m1.Load(key)
				if !ok || v != value {
					t.Fatal(v, ok, key, value)
				}
				return true
			})
		} else {
			v1, ok1 = m1.Load(r1)
			v2, ok2 = m2.Load(r1)
			if ok1 != ok2 || v1 != v2 {
				t.Fatal(rd, v1, ok1, v2, ok2)
			}
		}
	}
	// Correntness 3. (LoadOrStore)
	// Only one LoadorStore can successfully insert its key and value.
	// And the returned value is unique.
	mp := newset()
	tmpmap := newset()
	samekey := 123
	var added int64
	for i := 1; i < 1000; i++ {
		wg.Add(1)
		go func() {
			v := fastrand.Int63()
			actual, loaded := mp.LoadOrStore(samekey, v)
			if !loaded {
				atomic.AddInt64(&added, 1)
			}
			tmpmap.Store(int(actual.(int64)), nil)
			wg.Done()
		}()
	}
	wg.Wait()
	if added != 1 {
		t.Fatal("only one LoadOrStore can successfully insert a key and value")
	}
	if tmpmap.Len() != 1 {
		t.Fatal("only one value can be returned from LoadOrStore")
	}
	// Correntness 4. (LoadAndDelete)
	// Only one LoadAndDelete can successfully get a value.
	mp = newset()
	tmpmap = newset()
	samekey = 123
	added = 0 // int64
	mp.Store(samekey, 555)
	for i := 1; i < 1000; i++ {
		wg.Add(1)
		go func() {
			value, loaded := mp.LoadAndDelete(samekey)
			if loaded {
				atomic.AddInt64(&added, 1)
				if value != 555 {
					panic("invalid")
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
	if added != 1 {
		t.Fatal("Only one LoadAndDelete can successfully get a value")
	}

	// Correntness 5. (LoadOrStoreLazy)
	mp = newset()
	tmpmap = newset()
	samekey = 123
	added = 0
	var fcalled int64
	valuef := func() interface{} {
		atomic.AddInt64(&fcalled, 1)
		return fastrand.Int63()
	}
	for i := 1; i < 1000; i++ {
		wg.Add(1)
		go func() {
			actual, loaded := mp.LoadOrStoreLazy(samekey, valuef)
			if !loaded {
				atomic.AddInt64(&added, 1)
			}
			tmpmap.Store(int(actual.(int64)), nil)
			wg.Done()
		}()
	}
	wg.Wait()
	if added != 1 || fcalled != 1 {
		t.Fatal("only one LoadOrStoreLazy can successfully insert a key and value")
	}
	if tmpmap.Len() != 1 {
		t.Fatal("only one value can be returned from LoadOrStoreLazy")
	}
}

func testSkipMapIntDesc(t *testing.T, newset func() anyskipmap[int]) {
	m := newset()
	cases := []int{10, 11, 12}
	for _, v := range cases {
		m.Store(v, nil)
	}
	i := len(cases) - 1
	m.Range(func(key int, _ interface{}) bool {
		if key != cases[i] {
			t.Fail()
		}
		i--
		return true
	})
}

func testSkipMapString(t *testing.T, newset func() anyskipmap[string]) {
	m := newset()

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

/* Test from sync.Map */
func testSyncMapSuiteInt64(t *testing.T, newset func() anyskipmap[int64]) {
	const mapSize = 1 << 10

	m := newset()
	for n := int64(1); n <= mapSize; n++ {
		m.Store(n, int64(n))
	}

	done := make(chan struct{})
	var wg sync.WaitGroup
	defer func() {
		close(done)
		wg.Wait()
	}()
	for g := int64(runtime.GOMAXPROCS(0)); g > 0; g-- {
		r := rand.New(rand.NewSource(g))
		wg.Add(1)
		go func(g int64) {
			defer wg.Done()
			for i := int64(0); ; i++ {
				select {
				case <-done:
					return
				default:
				}
				for n := int64(1); n < mapSize; n++ {
					if r.Int63n(mapSize) == 0 {
						m.Store(n, n*i*g)
					} else {
						m.Load(n)
					}
				}
			}
		}(g)
	}

	iters := 1 << 10
	if testing.Short() {
		iters = 16
	}
	for n := iters; n > 0; n-- {
		seen := make(map[int64]bool, mapSize)

		m.Range(func(ki int64, vi interface{}) bool {
			k, v := ki, vi.(int64)
			if v%k != 0 {
				t.Fatalf("while Storing multiples of %v, Range saw value %v", k, v)
			}
			if seen[k] {
				t.Fatalf("Range visited key %v twice", k)
			}
			seen[k] = true
			return true
		})

		if len(seen) != mapSize {
			t.Fatalf("Range visited %v elements of %v-element Map", len(seen), mapSize)
		}
	}
}

func TestFloatMap(t *testing.T) {
	cases := []struct {
		k float64
		v int
	}{
		{math.NaN(), 1},
		{0.04, 1},
		{math.NaN(), -1},
		{0.05, 1},
		{math.Inf(1), 1},
		{0.04, 2},
		{math.NaN(), 2},
		{0.05, 2},
		{math.Inf(-1), -1},
		{math.Inf(1), 2},
		{math.Inf(-1), 2},
	}
	m := NewFloat64[int]()
	md := NewFloat64Desc[int]()
	m32 := NewFloat32[int]()
	m32d := NewFloat32Desc[int]()
	for _, kv := range cases {
		m.Store(kv.k, kv.v)
		md.Store(kv.k, kv.v)
		m32.Store(float32(kv.k), kv.v)
		m32d.Store(float32(kv.k), kv.v)
	}

	var (
		mr, mdr     []float64
		m32r, m32dr []float32
	)
	m.Range(func(key float64, value int) bool {
		mr = append(mr, key)
		if value != 2 {
			t.Fatal("invalid value", value)
		}
		return true
	})
	md.Range(func(key float64, value int) bool {
		mdr = append(mdr, key)
		if value != 2 {
			t.Fatal("invalid value", value)
		}
		return true
	})
	m32.Range(func(key float32, value int) bool {
		m32r = append(m32r, key)
		if value != 2 {
			t.Fatal("invalid value", value)
		}
		return true
	})
	m32d.Range(func(key float32, value int) bool {
		m32dr = append(m32dr, key)
		if value != 2 {
			t.Fatal("invalid value", value)
		}
		return true
	})

	var (
		asc = []float64{
			math.NaN(), math.Inf(-1), 0.04, 0.05, math.Inf(1),
		}
		desc = []float64{
			math.NaN(), math.Inf(1), 0.05, 0.04, math.Inf(-1),
		}
		asc32 = []float32{
			float32(math.NaN()), float32(math.Inf(-1)), 0.04, 0.05, float32(math.Inf(1)),
		}
		desc32 = []float32{
			float32(math.NaN()), float32(math.Inf(1)), 0.05, 0.04, float32(math.Inf(-1)),
		}
	)

	checkEqual := func(a, b []float64) {
		l := len(a)
		if len(b) != l {
			t.Fatal("invalid length", l)
		}
		for i := 0; i < l; i++ {
			if a[i] != b[i] && !(math.IsNaN(a[i])) {
				t.Fatal("not equal", i, a[i], b[i])
			}
		}
	}
	checkEqual32 := func(a, b []float32) {
		l := len(a)
		if len(b) != l {
			t.Fatal("invalid length", l)
		}
		for i := 0; i < l; i++ {
			if a[i] != b[i] && !(isNaNf32(a[i])) {
				t.Fatal("not equal", i, a[i], b[i])
			}
		}
	}
	checkEqual(mr, asc)
	checkEqual(mdr, desc)
	checkEqual32(m32r, asc32)
	checkEqual32(m32dr, desc32)
}
