// Code generated by gen.go; DO NOT EDIT.

package skipmap

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

// OrderedMap represents a map based on skip list.
type OrderedMap[T ordered] struct {
	length       int64
	highestLevel uint64 // highest level for now
	header       *orderednode[T]
}

type orderednode[T ordered] struct {
	value unsafe.Pointer // *any
	flags bitflag
	key   T
	next  optionalArray // [level]*orderednode
	mu    sync.Mutex
	level uint32
}

func newOrderedNode[T ordered](key T, value any, level int) *orderednode[T] {
	node := &orderednode[T]{
		key:   key,
		level: uint32(level),
	}
	node.storeVal(value)
	if level > op1 {
		node.next.extra = new([op2]unsafe.Pointer)
	}
	return node
}

func (n *orderednode[T]) storeVal(value any) {
	atomic.StorePointer(&n.value, unsafe.Pointer(&value))
}

func (n *orderednode[T]) loadVal() any {
	return *(*any)(atomic.LoadPointer(&n.value))
}

func (n *orderednode[T]) loadNext(i int) *orderednode[T] {
	return (*orderednode[T])(n.next.load(i))
}

func (n *orderednode[T]) storeNext(i int, node *orderednode[T]) {
	n.next.store(i, unsafe.Pointer(node))
}

func (n *orderednode[T]) atomicLoadNext(i int) *orderednode[T] {
	return (*orderednode[T])(n.next.atomicLoad(i))
}

func (n *orderednode[T]) atomicStoreNext(i int, node *orderednode[T]) {
	n.next.atomicStore(i, unsafe.Pointer(node))
}

// findNode takes a key and two maximal-height arrays then searches exactly as in a sequential skipmap.
// The returned preds and succs always satisfy preds[i] > key >= succs[i].
// (without fullpath, if find the node will return immediately)
func (s *OrderedMap[T]) findNode(key T, preds *[maxLevel]*orderednode[T], succs *[maxLevel]*orderednode[T]) *orderednode[T] {
	x := s.header
	for i := int(atomic.LoadUint64(&s.highestLevel)) - 1; i >= 0; i-- {
		succ := x.atomicLoadNext(i)
		for succ != nil && (succ.key < key) {
			x = succ
			succ = x.atomicLoadNext(i)
		}
		preds[i] = x
		succs[i] = succ

		// Check if the key already in the skipmap.
		if succ != nil && succ.key == key {
			return succ
		}
	}
	return nil
}

// findNodeDelete takes a key and two maximal-height arrays then searches exactly as in a sequential skip-list.
// The returned preds and succs always satisfy preds[i] > key >= succs[i].
func (s *OrderedMap[T]) findNodeDelete(key T, preds *[maxLevel]*orderednode[T], succs *[maxLevel]*orderednode[T]) int {
	// lFound represents the index of the first layer at which it found a node.
	lFound, x := -1, s.header
	for i := int(atomic.LoadUint64(&s.highestLevel)) - 1; i >= 0; i-- {
		succ := x.atomicLoadNext(i)
		for succ != nil && (succ.key < key) {
			x = succ
			succ = x.atomicLoadNext(i)
		}
		preds[i] = x
		succs[i] = succ

		// Check if the key already in the skip list.
		if lFound == -1 && succ != nil && succ.key == key {
			lFound = i
		}
	}
	return lFound
}

func unlockordered[T ordered](preds [maxLevel]*orderednode[T], highestLevel int) {
	var prevPred *orderednode[T]
	for i := highestLevel; i >= 0; i-- {
		if preds[i] != prevPred { // the node could be unlocked by previous loop
			preds[i].mu.Unlock()
			prevPred = preds[i]
		}
	}
}

// Store sets the value for a key.
func (s *OrderedMap[T]) Store(key T, value any) {
	level := s.randomlevel()
	var preds, succs [maxLevel]*orderednode[T]
	for {
		nodeFound := s.findNode(key, &preds, &succs)
		if nodeFound != nil { // indicating the key is already in the skip-list
			if !nodeFound.flags.Get(marked) {
				// We don't need to care about whether or not the node is fully linked,
				// just replace the value.
				nodeFound.storeVal(value)
				return
			}
			// If the node is marked, represents some other goroutines is in the process of deleting this node,
			// we need to add this node in next loop.
			continue
		}

		// Add this node into skip list.
		var (
			highestLocked        = -1 // the highest level being locked by this process
			valid                = true
			pred, succ, prevPred *orderednode[T]
		)
		for layer := 0; valid && layer < level; layer++ {
			pred = preds[layer]   // target node's previous node
			succ = succs[layer]   // target node's next node
			if pred != prevPred { // the node in this layer could be locked by previous loop
				pred.mu.Lock()
				highestLocked = layer
				prevPred = pred
			}
			// valid check if there is another node has inserted into the skip list in this layer during this process.
			// It is valid if:
			// 1. The previous node and next node both are not marked.
			// 2. The previous node's next node is succ in this layer.
			valid = !pred.flags.Get(marked) && (succ == nil || !succ.flags.Get(marked)) && pred.loadNext(layer) == succ
		}
		if !valid {
			unlockordered(preds, highestLocked)
			continue
		}

		nn := newOrderedNode(key, value, level)
		for layer := 0; layer < level; layer++ {
			nn.storeNext(layer, succs[layer])
			preds[layer].atomicStoreNext(layer, nn)
		}
		nn.flags.SetTrue(fullyLinked)
		unlockordered(preds, highestLocked)
		atomic.AddInt64(&s.length, 1)
	}
}

func (s *OrderedMap[T]) randomlevel() int {
	// Generate random level.
	level := randomLevel()
	// Update highest level if possible.
	for {
		hl := atomic.LoadUint64(&s.highestLevel)
		if uint64(level) <= hl {
			break
		}
		if atomic.CompareAndSwapUint64(&s.highestLevel, hl, uint64(level)) {
			break
		}
	}
	return level
}

// Load returns the value stored in the map for a key, or nil if no
// value is present.
// The ok result indicates whether value was found in the map.
func (s *OrderedMap[T]) Load(key T) (value any, ok bool) {
	x := s.header
	for i := int(atomic.LoadUint64(&s.highestLevel)) - 1; i >= 0; i-- {
		nex := x.atomicLoadNext(i)
		for nex != nil && (nex.key < key) {
			x = nex
			nex = x.atomicLoadNext(i)
		}

		// Check if the key already in the skip list.
		if nex != nil && nex.key == key {
			if nex.flags.MGet(fullyLinked|marked, fullyLinked) {
				return nex.loadVal(), true
			}
			return nil, false
		}
	}
	return nil, false
}

// LoadAndDelete deletes the value for a key, returning the previous value if any.
// The loaded result reports whether the key was present.
// (Modified from Delete)
func (s *OrderedMap[T]) LoadAndDelete(key T) (value any, loaded bool) {
	var (
		nodeToDelete *orderednode[T]
		isMarked     bool // represents if this operation mark the node
		topLayer     = -1
		preds, succs [maxLevel]*orderednode[T]
	)
	for {
		lFound := s.findNodeDelete(key, &preds, &succs)
		if isMarked || // this process mark this node or we can find this node in the skip list
			lFound != -1 && succs[lFound].flags.MGet(fullyLinked|marked, fullyLinked) && (int(succs[lFound].level)-1) == lFound {
			if !isMarked { // we don't mark this node for now
				nodeToDelete = succs[lFound]
				topLayer = lFound
				nodeToDelete.mu.Lock()
				if nodeToDelete.flags.Get(marked) {
					// The node is marked by another process,
					// the physical deletion will be accomplished by another process.
					nodeToDelete.mu.Unlock()
					return nil, false
				}
				nodeToDelete.flags.SetTrue(marked)
				isMarked = true
			}
			// Accomplish the physical deletion.
			var (
				highestLocked        = -1 // the highest level being locked by this process
				valid                = true
				pred, succ, prevPred *orderednode[T]
			)
			for layer := 0; valid && (layer <= topLayer); layer++ {
				pred, succ = preds[layer], succs[layer]
				if pred != prevPred { // the node in this layer could be locked by previous loop
					pred.mu.Lock()
					highestLocked = layer
					prevPred = pred
				}
				// valid check if there is another node has inserted into the skip list in this layer
				// during this process, or the previous is deleted by another process.
				// It is valid if:
				// 1. the previous node exists.
				// 2. no another node has inserted into the skip list in this layer.
				valid = !pred.flags.Get(marked) && pred.loadNext(layer) == succ
			}
			if !valid {
				unlockordered(preds, highestLocked)
				continue
			}
			for i := topLayer; i >= 0; i-- {
				// Now we own the `nodeToDelete`, no other goroutine will modify it.
				// So we don't need `nodeToDelete.loadNext`
				preds[i].atomicStoreNext(i, nodeToDelete.loadNext(i))
			}
			nodeToDelete.mu.Unlock()
			unlockordered(preds, highestLocked)
			atomic.AddInt64(&s.length, -1)
			return nodeToDelete.loadVal(), true
		}
		return nil, false
	}
}

// LoadOrStore returns the existing value for the key if present.
// Otherwise, it stores and returns the given value.
// The loaded result is true if the value was loaded, false if stored.
// (Modified from Store)
func (s *OrderedMap[T]) LoadOrStore(key T, value any) (actual any, loaded bool) {
	level := s.randomlevel()
	var preds, succs [maxLevel]*orderednode[T]
	for {
		nodeFound := s.findNode(key, &preds, &succs)
		if nodeFound != nil { // indicating the key is already in the skip-list
			if !nodeFound.flags.Get(marked) {
				// We don't need to care about whether or not the node is fully linked,
				// just return the value.
				return nodeFound.loadVal(), true
			}
			// If the node is marked, represents some other goroutines is in the process of deleting this node,
			// we need to add this node in next loop.
			continue
		}

		// Add this node into skip list.
		var (
			highestLocked        = -1 // the highest level being locked by this process
			valid                = true
			pred, succ, prevPred *orderednode[T]
		)
		for layer := 0; valid && layer < level; layer++ {
			pred = preds[layer]   // target node's previous node
			succ = succs[layer]   // target node's next node
			if pred != prevPred { // the node in this layer could be locked by previous loop
				pred.mu.Lock()
				highestLocked = layer
				prevPred = pred
			}
			// valid check if there is another node has inserted into the skip list in this layer during this process.
			// It is valid if:
			// 1. The previous node and next node both are not marked.
			// 2. The previous node's next node is succ in this layer.
			valid = !pred.flags.Get(marked) && (succ == nil || !succ.flags.Get(marked)) && pred.loadNext(layer) == succ
		}
		if !valid {
			unlockordered(preds, highestLocked)
			continue
		}

		nn := newOrderedNode(key, value, level)
		for layer := 0; layer < level; layer++ {
			nn.storeNext(layer, succs[layer])
			preds[layer].atomicStoreNext(layer, nn)
		}
		nn.flags.SetTrue(fullyLinked)
		unlockordered(preds, highestLocked)
		atomic.AddInt64(&s.length, 1)
		return value, false
	}
}

// LoadOrStoreLazy returns the existing value for the key if present.
// Otherwise, it stores and returns the given value from f, f will only be called once.
// The loaded result is true if the value was loaded, false if stored.
// (Modified from LoadOrStore)
func (s *OrderedMap[T]) LoadOrStoreLazy(key T, f func() any) (actual any, loaded bool) {
	level := s.randomlevel()
	var preds, succs [maxLevel]*orderednode[T]
	for {
		nodeFound := s.findNode(key, &preds, &succs)
		if nodeFound != nil { // indicating the key is already in the skip-list
			if !nodeFound.flags.Get(marked) {
				// We don't need to care about whether or not the node is fully linked,
				// just return the value.
				return nodeFound.loadVal(), true
			}
			// If the node is marked, represents some other goroutines is in the process of deleting this node,
			// we need to add this node in next loop.
			continue
		}

		// Add this node into skip list.
		var (
			highestLocked        = -1 // the highest level being locked by this process
			valid                = true
			pred, succ, prevPred *orderednode[T]
		)
		for layer := 0; valid && layer < level; layer++ {
			pred = preds[layer]   // target node's previous node
			succ = succs[layer]   // target node's next node
			if pred != prevPred { // the node in this layer could be locked by previous loop
				pred.mu.Lock()
				highestLocked = layer
				prevPred = pred
			}
			// valid check if there is another node has inserted into the skip list in this layer during this process.
			// It is valid if:
			// 1. The previous node and next node both are not marked.
			// 2. The previous node's next node is succ in this layer.
			valid = !pred.flags.Get(marked) && pred.loadNext(layer) == succ && (succ == nil || !succ.flags.Get(marked))
		}
		if !valid {
			unlockordered(preds, highestLocked)
			continue
		}
		value := f()
		nn := newOrderedNode(key, value, level)
		for layer := 0; layer < level; layer++ {
			nn.storeNext(layer, succs[layer])
			preds[layer].atomicStoreNext(layer, nn)
		}
		nn.flags.SetTrue(fullyLinked)
		unlockordered(preds, highestLocked)
		atomic.AddInt64(&s.length, 1)
		return value, false
	}
}

// Delete deletes the value for a key.
func (s *OrderedMap[T]) Delete(key T) bool {
	var (
		nodeToDelete *orderednode[T]
		isMarked     bool // represents if this operation mark the node
		topLayer     = -1
		preds, succs [maxLevel]*orderednode[T]
	)
	for {
		lFound := s.findNodeDelete(key, &preds, &succs)
		if isMarked || // this process mark this node or we can find this node in the skip list
			lFound != -1 && succs[lFound].flags.MGet(fullyLinked|marked, fullyLinked) && (int(succs[lFound].level)-1) == lFound {
			if !isMarked { // we don't mark this node for now
				nodeToDelete = succs[lFound]
				topLayer = lFound
				nodeToDelete.mu.Lock()
				if nodeToDelete.flags.Get(marked) {
					// The node is marked by another process,
					// the physical deletion will be accomplished by another process.
					nodeToDelete.mu.Unlock()
					return false
				}
				nodeToDelete.flags.SetTrue(marked)
				isMarked = true
			}
			// Accomplish the physical deletion.
			var (
				highestLocked        = -1 // the highest level being locked by this process
				valid                = true
				pred, succ, prevPred *orderednode[T]
			)
			for layer := 0; valid && (layer <= topLayer); layer++ {
				pred, succ = preds[layer], succs[layer]
				if pred != prevPred { // the node in this layer could be locked by previous loop
					pred.mu.Lock()
					highestLocked = layer
					prevPred = pred
				}
				// valid check if there is another node has inserted into the skip list in this layer
				// during this process, or the previous is deleted by another process.
				// It is valid if:
				// 1. the previous node exists.
				// 2. no another node has inserted into the skip list in this layer.
				valid = !pred.flags.Get(marked) && pred.atomicLoadNext(layer) == succ
			}
			if !valid {
				unlockordered(preds, highestLocked)
				continue
			}
			for i := topLayer; i >= 0; i-- {
				// Now we own the `nodeToDelete`, no other goroutine will modify it.
				// So we don't need `nodeToDelete.loadNext`
				preds[i].atomicStoreNext(i, nodeToDelete.loadNext(i))
			}
			nodeToDelete.mu.Unlock()
			unlockordered(preds, highestLocked)
			atomic.AddInt64(&s.length, -1)
			return true
		}
		return false
	}
}

// Range calls f sequentially for each key and value present in the skipmap.
// If f returns false, range stops the iteration.
//
// Range does not necessarily correspond to any consistent snapshot of the Map's
// contents: no key will be visited more than once, but if the value for any key
// is stored or deleted concurrently, Range may reflect any mapping for that key
// from any point during the Range call.
func (s *OrderedMap[T]) Range(f func(key T, value any) bool) {
	x := s.header.atomicLoadNext(0)
	for x != nil {
		if !x.flags.MGet(fullyLinked|marked, fullyLinked) {
			x = x.atomicLoadNext(0)
			continue
		}
		if !f(x.key, x.loadVal()) {
			break
		}
		x = x.atomicLoadNext(0)
	}
}

// Len returns the length of this skipmap.
func (s *OrderedMap[T]) Len() int {
	return int(atomic.LoadInt64(&s.length))
}