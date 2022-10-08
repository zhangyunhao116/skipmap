// Package skipmap is a high-performance, scalable, concurrent-safe map based on skip-list.
// In the typical pattern(100000 operations, 90%LOAD 9%STORE 1%DELETE, 8C16T), the skipmap
// up to 10x faster than the built-in sync.Map.
//
//go:generate go run gen.go
package skipmap

import "math"

// NewFunc returns an empty skipmap in ascending order.
//
// Note that the less function requires a strict weak ordering,
// see https://en.wikipedia.org/wiki/Weak_ordering#Strict_weak_orderings,
// or undefined behavior will happen.
func NewFunc[keyT any, valueT any](less func(a, b keyT) bool) *FuncMap[keyT, valueT] {
	var (
		t1 keyT
		t2 valueT
	)
	h := newFuncNode(t1, t2, maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &FuncMap[keyT, valueT]{
		header:       h,
		highestLevel: defaultHighestLevel,
		less:         less,
	}
}

// New returns an empty skipmap in ascending order.
func New[keyT ordered, valueT any]() *OrderedMap[keyT, valueT] {
	var (
		t1 keyT
		t2 valueT
	)
	h := newOrderedNode(t1, t2, maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &OrderedMap[keyT, valueT]{
		header:       h,
		highestLevel: defaultHighestLevel,
	}
}

// NewDesc returns an empty skipmap in descending order.
func NewDesc[keyT ordered, valueT any]() *OrderedMapDesc[keyT, valueT] {
	var (
		t1 keyT
		t2 valueT
	)
	h := newOrderedNodeDesc(t1, t2, maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &OrderedMapDesc[keyT, valueT]{
		header:       h,
		highestLevel: defaultHighestLevel,
	}
}

// NewString returns an empty skipmap in ascending order.
func NewString[valueT any]() *StringMap[valueT] {
	var t valueT
	h := newStringNode("", t, maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &StringMap[valueT]{
		header:       h,
		highestLevel: defaultHighestLevel,
	}
}

// NewStringDesc returns an empty skipmap in descending order.
func NewStringDesc[valueT any]() *StringMapDesc[valueT] {
	var t valueT
	h := newStringNodeDesc("", t, maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &StringMapDesc[valueT]{
		header:       h,
		highestLevel: defaultHighestLevel,
	}
}

func isNaNf32(x float32) bool {
	return x != x
}

// NewFloat32 returns an empty skipmap in ascending order.
func NewFloat32[valueT any]() *FuncMap[float32, valueT] {
	return NewFunc[float32, valueT](func(a, b float32) bool {
		return a < b || (isNaNf32(a) && !isNaNf32(b))
	})
}

// NewFloat32Desc returns an empty skipmap in descending order.
func NewFloat32Desc[valueT any]() *FuncMap[float32, valueT] {
	return NewFunc[float32, valueT](func(a, b float32) bool {
		return a > b || (isNaNf32(a) && !isNaNf32(b))
	})
}

// NewFloat64 returns an empty skipmap in ascending order.
func NewFloat64[valueT any]() *FuncMap[float64, valueT] {
	return NewFunc[float64, valueT](func(a, b float64) bool {
		return a < b || (math.IsNaN(a) && !math.IsNaN(b))
	})
}

// NewFloat64Desc returns an empty skipmap in descending order.
func NewFloat64Desc[valueT any]() *FuncMap[float64, valueT] {
	return NewFunc[float64, valueT](func(a, b float64) bool {
		return a > b || (math.IsNaN(a) && !math.IsNaN(b))
	})
}

// NewInt returns an empty skipmap in ascending order.
func NewInt[valueT any]() *IntMap[valueT] {
	var t valueT
	h := newIntNode(0, t, maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &IntMap[valueT]{
		header:       h,
		highestLevel: defaultHighestLevel,
	}
}

// NewIntDesc returns an empty skipmap in descending order.
func NewIntDesc[valueT any]() *IntMapDesc[valueT] {
	var t valueT
	h := newIntNodeDesc(0, t, maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &IntMapDesc[valueT]{
		header:       h,
		highestLevel: defaultHighestLevel,
	}
}

// NewInt64 returns an empty skipmap in ascending order.
func NewInt64[valueT any]() *Int64Map[valueT] {
	var t valueT
	h := newInt64Node(0, t, maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &Int64Map[valueT]{
		header:       h,
		highestLevel: defaultHighestLevel,
	}
}

// NewInt64Desc returns an empty skipmap in descending order.
func NewInt64Desc[valueT any]() *Int64MapDesc[valueT] {
	var t valueT
	h := newInt64NodeDesc(0, t, maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &Int64MapDesc[valueT]{
		header:       h,
		highestLevel: defaultHighestLevel,
	}
}

// NewInt32 returns an empty skipmap in ascending order.
func NewInt32[valueT any]() *Int32Map[valueT] {
	var t valueT
	h := newInt32Node(0, t, maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &Int32Map[valueT]{
		header:       h,
		highestLevel: defaultHighestLevel,
	}
}

// NewInt32Desc returns an empty skipmap in descending order.
func NewInt32Desc[valueT any]() *Int32MapDesc[valueT] {
	var t valueT
	h := newInt32NodeDesc(0, t, maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &Int32MapDesc[valueT]{
		header:       h,
		highestLevel: defaultHighestLevel,
	}
}

// NewUint64 returns an empty skipmap in ascending order.
func NewUint64[valueT any]() *Uint64Map[valueT] {
	var t valueT
	h := newUint64Node(0, t, maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &Uint64Map[valueT]{
		header:       h,
		highestLevel: defaultHighestLevel,
	}
}

// NewUint64Desc returns an empty skipmap in descending order.
func NewUint64Desc[valueT any]() *Uint64MapDesc[valueT] {
	var t valueT
	h := newUint64NodeDesc(0, t, maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &Uint64MapDesc[valueT]{
		header:       h,
		highestLevel: defaultHighestLevel,
	}
}

// NewUint32 returns an empty skipmap in ascending order.
func NewUint32[valueT any]() *Uint32Map[valueT] {
	var t valueT
	h := newUint32Node(0, t, maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &Uint32Map[valueT]{
		header:       h,
		highestLevel: defaultHighestLevel,
	}
}

// NewUint32Desc returns an empty skipmap in descending order.
func NewUint32Desc[valueT any]() *Uint32MapDesc[valueT] {
	var t valueT
	h := newUint32NodeDesc(0, t, maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &Uint32MapDesc[valueT]{
		header:       h,
		highestLevel: defaultHighestLevel,
	}
}

// NewUint returns an empty skipmap in ascending order.
func NewUint[valueT any]() *UintMap[valueT] {
	var t valueT
	h := newUintNode(0, t, maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &UintMap[valueT]{
		header:       h,
		highestLevel: defaultHighestLevel,
	}
}

// NewUintDesc returns an empty skipmap in descending order.
func NewUintDesc[valueT any]() *UintMapDesc[valueT] {
	var t valueT
	h := newUintNodeDesc(0, t, maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &UintMapDesc[valueT]{
		header:       h,
		highestLevel: defaultHighestLevel,
	}
}
