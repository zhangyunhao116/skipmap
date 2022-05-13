// Package skipmap is a high-performance, scalable, concurrent-safe map based on skip-list.
// In the typical pattern(100000 operations, 90%LOAD 9%STORE 1%DELETE, 8C16T), the skipmap
// up to 10x faster than the built-in sync.Map.
package skipmap

// NewFunc returns an empty skipmap in ascending order.
//
// Note that the less function requires a strict weak ordering,
// see https://en.wikipedia.org/wiki/Weak_ordering#Strict_weak_orderings,
// or undefined behavior will happen.
func NewFunc[T any](less func(a, b T) bool) *FuncMap[T] {
	var t T
	h := newFuncNode(t, nil, maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &FuncMap[T]{
		header:       h,
		highestLevel: defaultHighestLevel,
		less:         less,
	}
}

// New returns an empty skipmap in ascending order.
func New[T ordered]() *OrderedMap[T] {
	var t T
	h := newOrderedNode(t, nil, maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &OrderedMap[T]{
		header:       h,
		highestLevel: defaultHighestLevel,
	}
}

// NewDesc returns an empty skipmap in descending order.
func NewDesc[T ordered]() *OrderedMapDesc[T] {
	var t T
	h := newOrderedNodeDesc(t, nil, maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &OrderedMapDesc[T]{
		header:       h,
		highestLevel: defaultHighestLevel,
	}
}

// NewString returns an empty skipmap in ascending order.
func NewString() *StringMap {
	h := newStringNode("", nil, maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &StringMap{
		header:       h,
		highestLevel: defaultHighestLevel,
	}
}

// NewString returns an empty skipmap in descending order.
func NewStringDesc() *StringMapDesc {
	h := newStringNodeDesc("", nil, maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &StringMapDesc{
		header:       h,
		highestLevel: defaultHighestLevel,
	}
}

// NewFloat32 returns an empty skipmap in ascending order.
func NewFloat32() *Float32Map {
	h := newFloat32Node(0, nil, maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &Float32Map{
		header:       h,
		highestLevel: defaultHighestLevel,
	}
}

// NewFloat32 returns an empty skipmap in descending order.
func NewFloat32Desc() *Float32MapDesc {
	h := newFloat32NodeDesc(0, nil, maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &Float32MapDesc{
		header:       h,
		highestLevel: defaultHighestLevel,
	}
}

// NewFloat64 returns an empty skipmap in ascending order.
func NewFloat64() *Float64Map {
	h := newFloat64Node(0, nil, maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &Float64Map{
		header:       h,
		highestLevel: defaultHighestLevel,
	}
}

// NewFloat64 returns an empty skipmap in descending order.
func NewFloat64Desc() *Float64MapDesc {
	h := newFloat64NodeDesc(0, nil, maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &Float64MapDesc{
		header:       h,
		highestLevel: defaultHighestLevel,
	}
}

// NewInt returns an empty skipmap in ascending order.
func NewInt() *IntMap {
	h := newIntNode(0, nil, maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &IntMap{
		header:       h,
		highestLevel: defaultHighestLevel,
	}
}

// NewInt returns an empty skipmap in descending order.
func NewIntDesc() *IntMapDesc {
	h := newIntNodeDesc(0, nil, maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &IntMapDesc{
		header:       h,
		highestLevel: defaultHighestLevel,
	}
}

// NewInt64 returns an empty skipmap in ascending order.
func NewInt64() *Int64Map {
	h := newInt64Node(0, nil, maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &Int64Map{
		header:       h,
		highestLevel: defaultHighestLevel,
	}
}

// NewInt64 returns an empty skipmap in descending order.
func NewInt64Desc() *Int64MapDesc {
	h := newInt64NodeDesc(0, nil, maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &Int64MapDesc{
		header:       h,
		highestLevel: defaultHighestLevel,
	}
}

// NewInt32 returns an empty skipmap in ascending order.
func NewInt32() *Int32Map {
	h := newInt32Node(0, nil, maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &Int32Map{
		header:       h,
		highestLevel: defaultHighestLevel,
	}
}

// NewInt32 returns an empty skipmap in descending order.
func NewInt32Desc() *Int32MapDesc {
	h := newInt32NodeDesc(0, nil, maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &Int32MapDesc{
		header:       h,
		highestLevel: defaultHighestLevel,
	}
}

// NewUint64 returns an empty skipmap in ascending order.
func NewUint64() *Uint64Map {
	h := newUint64Node(0, nil, maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &Uint64Map{
		header:       h,
		highestLevel: defaultHighestLevel,
	}
}

// NewUint64 returns an empty skipmap in descending order.
func NewUint64Desc() *Uint64MapDesc {
	h := newUint64NodeDesc(0, nil, maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &Uint64MapDesc{
		header:       h,
		highestLevel: defaultHighestLevel,
	}
}

// NewUint32 returns an empty skipmap in ascending order.
func NewUint32() *Uint32Map {
	h := newUint32Node(0, nil, maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &Uint32Map{
		header:       h,
		highestLevel: defaultHighestLevel,
	}
}

// NewUint32 returns an empty skipmap in descending order.
func NewUint32Desc() *Uint32MapDesc {
	h := newUint32NodeDesc(0, nil, maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &Uint32MapDesc{
		header:       h,
		highestLevel: defaultHighestLevel,
	}
}

// NewUint returns an empty skipmap in ascending order.
func NewUint() *UintMap {
	h := newUintNode(0, nil, maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &UintMap{
		header:       h,
		highestLevel: defaultHighestLevel,
	}
}

// NewUint returns an empty skipmap in descending order.
func NewUintDesc() *UintMapDesc {
	h := newUintNodeDesc(0, nil, maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &UintMapDesc{
		header:       h,
		highestLevel: defaultHighestLevel,
	}
}
