package skipmap

import (
	"math"
	"strconv"
	"sync"
	"testing"
)

const initsize = 1 << 10 // for `load` `1Delete9Store90Load` `1Range9Delete90Store900Load`
const randN = math.MaxUint32

func BenchmarkStore(b *testing.B) {
	b.Run("skipmap", func(b *testing.B) {
		l := NewInt64()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				l.Store(int64(fastrandn(randN)), nil)
			}
		})
	})
	b.Run("sync.Map", func(b *testing.B) {
		var l sync.Map
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				l.Store(int64(fastrandn(randN)), nil)
			}
		})
	})
}

func BenchmarkLoad100Hits(b *testing.B) {
	b.Run("skipmap", func(b *testing.B) {
		l := NewInt64()
		for i := 0; i < initsize; i++ {
			l.Store(int64(i), nil)
		}
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, _ = l.Load(int64(fastrandn(initsize)))
			}
		})
	})
	b.Run("sync.Map", func(b *testing.B) {
		var l sync.Map
		for i := 0; i < initsize; i++ {
			l.Store(int64(i), nil)
		}
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, _ = l.Load(int64(fastrandn(initsize)))
			}
		})
	})
}

func BenchmarkLoad50Hits(b *testing.B) {
	const rate = 2
	b.Run("skipmap", func(b *testing.B) {
		l := NewInt64()
		for i := 0; i < initsize*rate; i++ {
			if fastrandn(rate) == 0 {
				l.Store(int64(i), nil)
			}
		}
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, _ = l.Load(int64(fastrandn(initsize * rate)))
			}
		})
	})
	b.Run("sync.Map", func(b *testing.B) {
		var l sync.Map
		for i := 0; i < initsize*rate; i++ {
			if fastrandn(rate) == 0 {
				l.Store(int64(i), nil)
			}
		}
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, _ = l.Load(int64(fastrandn(initsize * rate)))
			}
		})
	})
}

func BenchmarkLoadNoHits(b *testing.B) {
	b.Run("skipmap", func(b *testing.B) {
		l := NewInt64()
		invalid := make([]int64, 0, initsize)
		for i := 0; i < initsize*2; i++ {
			if i%2 == 0 {
				l.Store(int64(i), nil)
			} else {
				invalid = append(invalid, int64(i))
			}
		}
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, _ = l.Load(invalid[fastrandn(uint32(len(invalid)))])
			}
		})
	})
	b.Run("sync.Map", func(b *testing.B) {
		var l sync.Map
		invalid := make([]int64, 0, initsize)
		for i := 0; i < initsize*2; i++ {
			if i%2 == 0 {
				l.Store(int64(i), nil)
			} else {
				invalid = append(invalid, int64(i))
			}
		}
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, _ = l.Load(invalid[fastrandn(uint32(len(invalid)))])
			}
		})
	})
}

func Benchmark50Store50Load(b *testing.B) {
	b.Run("skipmap", func(b *testing.B) {
		l := NewInt64()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				u := fastrandn(10)
				if u < 5 {
					l.Store(int64(fastrandn(randN)), nil)
				} else {
					l.Load(int64(fastrandn(randN)))
				}
			}
		})
	})
	b.Run("sync.Map", func(b *testing.B) {
		var l sync.Map
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				u := fastrandn(10)
				if u < 5 {
					l.Store(int64(fastrandn(randN)), nil)
				} else {
					l.Load(int64(fastrandn(randN)))
				}
			}
		})
	})
}

func Benchmark30Store70Load(b *testing.B) {
	b.Run("skipmap", func(b *testing.B) {
		l := NewInt64()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				u := fastrandn(10)
				if u < 3 {
					l.Store(int64(fastrandn(randN)), nil)
				} else {
					l.Load(int64(fastrandn(randN)))
				}
			}
		})
	})
	b.Run("sync.Map", func(b *testing.B) {
		var l sync.Map
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				u := fastrandn(10)
				if u < 3 {
					l.Store(int64(fastrandn(randN)), nil)
				} else {
					l.Load(int64(fastrandn(randN)))
				}
			}
		})
	})
}

func Benchmark1Delete9Store90Load(b *testing.B) {
	b.Run("skipmap", func(b *testing.B) {
		l := NewInt64()
		for i := 0; i < initsize; i++ {
			l.Store(int64(i), nil)
		}
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				u := fastrandn(100)
				if u < 9 {
					l.Store(int64(fastrandn(randN)), nil)
				} else if u == 2 {
					l.Delete(int64(fastrandn(randN)))
				} else {
					l.Load(int64(fastrandn(randN)))
				}
			}
		})
	})
	b.Run("sync.Map", func(b *testing.B) {
		var l sync.Map
		for i := 0; i < initsize; i++ {
			l.Store(int64(i), nil)
		}
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				u := fastrandn(100)
				if u < 9 {
					l.Store(int64(fastrandn(randN)), nil)
				} else if u == 2 {
					l.Delete(int64(fastrandn(randN)))
				} else {
					l.Load(int64(fastrandn(randN)))
				}
			}
		})
	})
}

func Benchmark1Range9Delete90Store900Load(b *testing.B) {
	b.Run("skipmap", func(b *testing.B) {
		l := NewInt64()
		for i := 0; i < initsize; i++ {
			l.Store(int64(i), nil)
		}
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				u := fastrandn(1000)
				if u == 0 {
					l.Range(func(key int64, value interface{}) bool {
						return true
					})
				} else if u > 10 && u < 20 {
					l.Delete(int64(fastrandn(randN)))
				} else if u >= 100 && u < 190 {
					l.Store(int64(fastrandn(randN)), nil)
				} else {
					l.Load(int64(fastrandn(randN)))
				}
			}
		})
	})
	b.Run("sync.Map", func(b *testing.B) {
		var l sync.Map
		for i := 0; i < initsize; i++ {
			l.Store(int64(i), nil)
		}
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				u := fastrandn(1000)
				if u == 0 {
					l.Range(func(key, value interface{}) bool {
						return true
					})
				} else if u > 10 && u < 20 {
					l.Delete(int64(fastrandn(randN)))
				} else if u >= 100 && u < 190 {
					l.Store(int64(fastrandn(randN)), nil)
				} else {
					l.Load(int64(fastrandn(randN)))
				}
			}
		})
	})
}

func BenchmarkStringStore(b *testing.B) {
	b.Run("skipmap", func(b *testing.B) {
		l := NewString()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				l.Store(strconv.Itoa(int(fastrand())), nil)
			}
		})
	})
	b.Run("sync.Map", func(b *testing.B) {
		var l sync.Map
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				l.Store(strconv.Itoa(int(fastrand())), nil)
			}
		})
	})
}

func BenchmarkStringLoad50Hits(b *testing.B) {
	const rate = 2
	b.Run("skipmap", func(b *testing.B) {
		l := NewString()
		for i := 0; i < initsize*rate; i++ {
			if fastrandn(rate) == 0 {
				l.Store(strconv.Itoa(i), nil)
			}
		}
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, _ = l.Load(strconv.Itoa(int(fastrandn(initsize * rate))))
			}
		})
	})
	b.Run("sync.Map", func(b *testing.B) {
		var l sync.Map
		for i := 0; i < initsize*rate; i++ {
			if fastrandn(rate) == 0 {
				l.Store(strconv.Itoa(i), nil)
			}
		}
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, _ = l.Load(strconv.Itoa(int(fastrandn(initsize * rate))))
			}
		})
	})
}

func BenchmarkString30Store70Load(b *testing.B) {
	b.Run("skipmap", func(b *testing.B) {
		l := NewString()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				u := fastrandn(10)
				if u < 3 {
					l.Store(strconv.Itoa(int(fastrandn(randN))), nil)
				} else {
					l.Load(strconv.Itoa(int(fastrandn(randN))))
				}
			}
		})
	})
	b.Run("sync.Map", func(b *testing.B) {
		var l sync.Map
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				u := fastrandn(10)
				if u < 3 {
					l.Store(strconv.Itoa(int(fastrandn(randN))), nil)
				} else {
					l.Load(strconv.Itoa(int(fastrandn(randN))))
				}
			}
		})
	})
}

func BenchmarkString1Delete9Store90Load(b *testing.B) {
	b.Run("skipmap", func(b *testing.B) {
		l := NewString()
		for i := 0; i < initsize; i++ {
			l.Store(strconv.Itoa(i), nil)
		}
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				u := fastrandn(100)
				if u == 1 {
					l.Store(strconv.Itoa(int(fastrandn(randN))), nil)
				} else if u == 2 {
					l.Delete(strconv.Itoa(int(fastrandn(randN))))
				} else {
					l.Load(strconv.Itoa(int(fastrandn(randN))))
				}
			}
		})
	})
	b.Run("sync.Map", func(b *testing.B) {
		var l sync.Map
		for i := 0; i < initsize; i++ {
			l.Store(strconv.Itoa(i), nil)
		}
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				u := fastrandn(100)
				if u == 1 {
					l.Store(strconv.Itoa(int(fastrandn(randN))), nil)
				} else if u == 2 {
					l.Delete(strconv.Itoa(int(fastrandn(randN))))
				} else {
					l.Load(strconv.Itoa(int(fastrandn(randN))))
				}
			}
		})
	})
}

func BenchmarkString1Range9Delete90Store900Load(b *testing.B) {
	b.Run("skipmap", func(b *testing.B) {
		l := NewString()
		for i := 0; i < initsize; i++ {
			l.Store(strconv.Itoa(i), nil)
		}
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				u := fastrandn(1000)
				if u == 0 {
					l.Range(func(key string, value interface{}) bool {
						return true
					})
				} else if u > 10 && u < 20 {
					l.Delete(strconv.Itoa(int(fastrandn(randN))))
				} else if u >= 100 && u < 190 {
					l.Store(strconv.Itoa(int(fastrandn(randN))), nil)
				} else {
					l.Load(strconv.Itoa(int(fastrandn(randN))))
				}
			}
		})
	})
	b.Run("sync.Map", func(b *testing.B) {
		var l sync.Map
		for i := 0; i < initsize; i++ {
			l.Store(strconv.Itoa(i), nil)
		}
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				u := fastrandn(1000)
				if u == 0 {
					l.Range(func(key, value interface{}) bool {
						return true
					})
				} else if u > 10 && u < 20 {
					l.Delete(strconv.Itoa(int(fastrandn(randN))))
				} else if u >= 100 && u < 190 {
					l.Store(strconv.Itoa(int(fastrandn(randN))), nil)
				} else {
					l.Load(strconv.Itoa(int(fastrandn(randN))))
				}
			}
		})
	})
}
