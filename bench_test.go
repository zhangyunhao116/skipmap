package skipmap

import (
	"testing"

	"github.com/zhangyunhao116/fastrand"
)

func BenchmarkLoadOrStoreExist(b *testing.B) {
	m := NewInt[int]()
	m.Store(1, 1)
	b.ResetTimer()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			m.LoadOrStore(1, 1)
		}
	})
}

func BenchmarkLoadOrStoreLazyExist(b *testing.B) {
	m := NewInt[int]()
	m.Store(1, 1)
	b.ResetTimer()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			m.LoadOrStoreLazy(1, func() int { return 1 })
		}
	})
}

func BenchmarkLoadOrStoreExistSingle(b *testing.B) {
	m := NewInt[int]()
	m.Store(1, 1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.LoadOrStore(1, 1)
	}
}

func BenchmarkLoadOrStoreLazyExistSingle(b *testing.B) {
	m := NewInt[int]()
	m.Store(1, 1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.LoadOrStoreLazy(1, func() int { return 1 })
	}
}

func BenchmarkLoadOrStoreRandom(b *testing.B) {
	m := NewInt[int]()
	b.ResetTimer()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			m.LoadOrStore(fastrand.Int(), 1)
		}
	})
}

func BenchmarkLoadOrStoreLazyRandom(b *testing.B) {
	m := NewInt[int]()
	b.ResetTimer()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			m.LoadOrStoreLazy(fastrand.Int(), func() int { return 1 })
		}
	})
}

func BenchmarkLoadOrStoreRandomSingle(b *testing.B) {
	m := NewInt[int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.LoadOrStore(fastrand.Int(), 1)
	}
}

func BenchmarkLoadOrStoreLazyRandomSingle(b *testing.B) {
	m := NewInt[int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.LoadOrStoreLazy(fastrand.Int(), func() int { return 1 })
	}
}
