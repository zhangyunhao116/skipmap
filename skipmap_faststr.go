package skipmap

import (
	"unsafe"

	"github.com/zhangyunhao116/fastrand"
	"github.com/zhangyunhao116/skipmap/internal/typehack"
)

var strhash func(string) uint64

func init() {
	runtimestrhash := typehack.NewHasher[string]()
	randomseed := fastrand.Uint()
	strhash = func(s string) uint64 {
		return uint64(runtimestrhash(unsafe.Pointer(&s), uintptr(randomseed)))
	}
}

// NewStringFast returns an empty skipmap with string key.
// The item order of the skipmap is different between each run.
// If you need to keep the item order of each run, use [`NewString`].
// The [`StringMapFast`] is about 25% faster than the [`StringMap`].
func NewStringFast[valueT any]() *StringMapFast[valueT] {
	var t valueT
	h := newStringNodeFast(0, "", t, maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &StringMapFast[valueT]{
		header:       h,
		highestLevel: defaultHighestLevel,
	}
}
