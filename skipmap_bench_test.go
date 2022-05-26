package skipmap

import (
	"math"
	"strconv"
	"sync"
	"testing"

	"github.com/zhangyunhao116/fastrand"
)

const (
	initsize         = 1 << 10 // for `Load` `1Delete9Store90Load` `1Range9Delete90Store900Load`
	randN            = math.MaxUint32
	mockvalue uint64 = 100
)

func BenchmarkInt64(b *testing.B) {
	all := []benchInt64Task{{
		name: "skipmap", New: func() int64Map {
			return NewInt64[any]()
		}}}
	all = append(all, benchInt64Task{
		name: "sync.Map", New: func() int64Map {
			return new(int64SyncMap)
		}})
	benchStore(b, all)
	benchLoad50Hits(b, all)
	bench30Store70Load(b, all)
	bench1Delete9Store90Load(b, all)
	bench1Range9Delete90Store900Load(b, all)
}

func BenchmarkString(b *testing.B) {
	all := []benchStringTask{{
		name: "skipmap", New: func() stringMap {
			return NewString[any]()
		}}}
	all = append(all, benchStringTask{
		name: "sync.Map", New: func() stringMap {
			return new(stringSyncMap)
		}})
	benchStringStore(b, all)
	benchStringLoad50Hits(b, all)
	benchString30Store70Load(b, all)
	benchString1Delete9Store90Load(b, all)
	benchString1Range9Delete90Store900Load(b, all)
}

func benchStore(b *testing.B, benchTasks []benchInt64Task) {
	for _, v := range benchTasks {
		b.Run("Store/"+v.name, func(b *testing.B) {
			s := v.New()
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					s.Store(int64(fastrand.Uint32n(randN)), mockvalue)
				}
			})
		})
	}
}

func benchLoad50Hits(b *testing.B, benchTasks []benchInt64Task) {
	for _, v := range benchTasks {
		b.Run("Load50Hits/"+v.name, func(b *testing.B) {
			const rate = 2
			s := v.New()
			for i := 0; i < initsize*rate; i++ {
				if fastrand.Uint32n(rate) == 0 {
					s.Store(int64(i), mockvalue)
				}
			}
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					s.Load(int64(fastrand.Uint32n(initsize * rate)))
				}
			})
		})
	}
}

func bench30Store70Load(b *testing.B, benchTasks []benchInt64Task) {
	for _, v := range benchTasks {
		b.Run("30Store70Load/"+v.name, func(b *testing.B) {
			s := v.New()
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					u := fastrand.Uint32n(10)
					if u < 3 {
						s.Store(int64(fastrand.Uint32n(randN)), mockvalue)
					} else {
						s.Load(int64(fastrand.Uint32n(randN)))
					}
				}
			})
		})
	}
}

func bench1Delete9Store90Load(b *testing.B, benchTasks []benchInt64Task) {
	for _, v := range benchTasks {
		b.Run("1Delete9Store90Load/"+v.name, func(b *testing.B) {
			s := v.New()
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					u := fastrand.Uint32n(100)
					if u < 9 {
						s.Store(int64(fastrand.Uint32n(randN)), mockvalue)
					} else if u == 10 {
						s.Delete(int64(fastrand.Uint32n(randN)))
					} else {
						s.Load(int64(fastrand.Uint32n(randN)))
					}
				}
			})
		})
	}
}

func bench1Range9Delete90Store900Load(b *testing.B, benchTasks []benchInt64Task) {
	for _, v := range benchTasks {
		b.Run("1Range9Delete90Store900Load/"+v.name, func(b *testing.B) {
			s := v.New()
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					u := fastrand.Uint32n(1000)
					if u == 0 {
						s.Range(func(key int64, value interface{}) bool {
							return true
						})
					} else if u > 10 && u < 20 {
						s.Delete(int64(fastrand.Uint32n(randN)))
					} else if u >= 100 && u < 190 {
						s.Store(int64(fastrand.Uint32n(randN)), mockvalue)
					} else {
						s.Load(int64(fastrand.Uint32n(randN)))
					}
				}
			})
		})
	}
}

func benchStringStore(b *testing.B, benchTasks []benchStringTask) {
	for _, v := range benchTasks {
		b.Run("Store/"+v.name, func(b *testing.B) {
			s := v.New()
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					s.Store(strconv.Itoa(int(fastrand.Uint32n(randN))), mockvalue)
				}
			})
		})
	}
}

func benchStringLoad50Hits(b *testing.B, benchTasks []benchStringTask) {
	for _, v := range benchTasks {
		b.Run("Load50Hits/"+v.name, func(b *testing.B) {
			const rate = 2
			s := v.New()
			for i := 0; i < initsize*rate; i++ {
				if fastrand.Uint32n(rate) == 0 {
					s.Store(strconv.Itoa(int(fastrand.Uint32n(randN))), mockvalue)
				}
			}
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					s.Load(strconv.Itoa(int(fastrand.Uint32n(randN))))
				}
			})
		})
	}
}

func benchString30Store70Load(b *testing.B, benchTasks []benchStringTask) {
	for _, v := range benchTasks {
		b.Run("30Store70Load/"+v.name, func(b *testing.B) {
			s := v.New()
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					u := fastrand.Uint32n(10)
					if u < 3 {
						s.Store(strconv.Itoa(int(fastrand.Uint32n(randN))), mockvalue)
					} else {
						s.Load(strconv.Itoa(int(fastrand.Uint32n(randN))))
					}
				}
			})
		})
	}
}

func benchString1Delete9Store90Load(b *testing.B, benchTasks []benchStringTask) {
	for _, v := range benchTasks {
		b.Run("1Delete9Store90Load/"+v.name, func(b *testing.B) {
			s := v.New()
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					u := fastrand.Uint32n(100)
					if u < 9 {
						s.Store(strconv.Itoa(int(fastrand.Uint32n(randN))), mockvalue)
					} else if u == 10 {
						s.Delete(strconv.Itoa(int(fastrand.Uint32n(randN))))
					} else {
						s.Load(strconv.Itoa(int(fastrand.Uint32n(randN))))
					}
				}
			})
		})
	}
}

func benchString1Range9Delete90Store900Load(b *testing.B, benchTasks []benchStringTask) {
	for _, v := range benchTasks {
		b.Run("1Range9Delete90Store900Load/"+v.name, func(b *testing.B) {
			s := v.New()
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					u := fastrand.Uint32n(1000)
					if u == 0 {
						s.Range(func(_ string, _ interface{}) bool {
							return true
						})
					} else if u > 10 && u < 20 {
						s.Delete(strconv.Itoa(int(fastrand.Uint32n(randN))))
					} else if u >= 100 && u < 190 {
						s.Store(strconv.Itoa(int(fastrand.Uint32n(randN))), mockvalue)
					} else {
						s.Load(strconv.Itoa(int(fastrand.Uint32n(randN))))
					}
				}
			})
		})
	}
}

type benchInt64Task struct {
	name string
	New  func() int64Map
}

type int64Map interface {
	Store(x int64, v interface{})
	Load(x int64) (interface{}, bool)
	Delete(x int64) bool
	Range(f func(key int64, value interface{}) bool)
}

type int64SyncMap struct {
	data sync.Map
}

func (m *int64SyncMap) Store(x int64, v interface{}) {
	m.data.Store(x, v)
}

func (m *int64SyncMap) Load(x int64) (interface{}, bool) {
	return m.data.Load(x)
}

func (m *int64SyncMap) Delete(x int64) bool {
	m.data.Delete(x)
	return true
}

func (m *int64SyncMap) Range(f func(key int64, value interface{}) bool) {
	m.data.Range(func(key, value interface{}) bool {
		return !f(key.(int64), value)
	})
}

type benchStringTask struct {
	name string
	New  func() stringMap
}

type stringMap interface {
	Store(x string, v interface{})
	Load(x string) (interface{}, bool)
	Delete(x string) bool
	Range(f func(key string, value interface{}) bool)
}

type stringSyncMap struct {
	data sync.Map
}

func (m *stringSyncMap) Store(x string, v interface{}) {
	m.data.Store(x, v)
}

func (m *stringSyncMap) Load(x string) (interface{}, bool) {
	return m.data.Load(x)
}

func (m *stringSyncMap) Delete(x string) bool {
	m.data.Delete(x)
	return true
}

func (m *stringSyncMap) Range(f func(key string, value interface{}) bool) {
	m.data.Range(func(key, value interface{}) bool {
		return !f(key.(string), value)
	})
}
