<p align="center">
  <img src="https://raw.githubusercontent.com/zhangyunhao116/public-data/master/skipmap-logo.png"/>
</p>

## Introduction

> If your Go version is lower than 1.18, use v0.7.0 instead.

skipmap is a high-performance, scalable, concurrent-safe map based on skip-list. In the typical pattern(100000 operations, 90%LOAD 9%STORE 1%DELETE, 8C16T), the skipmap up to 10x faster than the built-in sync.Map.

The main idea behind the skipmap is [A Simple Optimistic Skiplist Algorithm](<https://people.csail.mit.edu/shanir/publications/LazySkipList.pdf>).

Different from the sync.Map, the keys in the skipmap are always sorted, and the `Load` and `Range` operations are wait-free (A goroutine is guaranteed to complete an operation as long as it keeps taking steps, regardless of the activity of other goroutines).


## Features

- Scalable, high-performance, concurrent-safe.
- Wait-free Load and Range operations (wait-free algorithms have stronger guarantees than lock-free).
- Sorted items.



## When should you use skipmap

In most cases, `skipmap` is better than `sync.Map`, especially in these situations: 

- **Sorted keys are needed**.
- **Concurrent calls multiple operations**. Such as use both `Range` and `Store` at the same time, in this situation, use skipmap can obtain very large improvement on performance.

If only one goroutine access the map for the most of the time, such as insert a batch of elements and then use only `Load` or `Range`, use built-in map is better.


## QuickStart

See [Go doc](https://pkg.go.dev/github.com/zhangyunhao116/skipmap) for more information.

```go
package main

import (
	"fmt"

	"github.com/zhangyunhao116/skipmap"
)

func main() {
	// Typed key and generic value.
	m0 := skipmap.NewString[int]()

	for _, v := range []int{10, 12, 15} {
		m0.Store(strconv.Itoa(v), v+100)
	}

	v, ok := m0.Load("10")
	if ok {
		fmt.Println("skipmap load 10 with value ", v)
	}

	m0.Range(func(key string, value int) bool {
		fmt.Println("skipmap range found ", key, value)
		return true
	})

	m0.Delete("15")
	fmt.Printf("skipmap contains %d items\r\n", m0.Len())

	// Generic key and value.
	m1 := skipmap.New[string, int]()
	for _, v := range []int{11, 13, 16} {
		m1.Store(strconv.Itoa(v), v+100)
	}
	m1.Range(func(key string, value int) bool {
		println("m1 found ", key, value)
		return true
	})

	// Generic key and value with less function.
	m2 := skipmap.NewFunc[int, string](func(a, b int) bool { return a < b })
	for _, v := range []int{15, 17, 19} {
		m2.Store(v, strconv.Itoa(v+200))
	}
	m2.Range(func(key int, value string) bool {
		println("m2 found ", key, value)
		return true
	})
}

```

**Note that generic APIs are always slower than typed APIs, but are more suitable for some scenarios such as functional programming.**

> e.g. `New[string,int]` is \~2x slower than `NewString[int]`, and `NewFunc[string,int](func(a, b string) bool { return a < b })` is 1\~2x slower than `NewString[int]`.
>
> Performance ranking: `NewString[int]` > `New[string,int]` > `NewFunc[string,int](func(a, b string) bool { return a < b })`


## Benchmark

> based on typed APIs.

Go version: go1.16.2 linux/amd64

CPU: AMD 3700x(8C16T), running at 3.6GHz

OS: ubuntu 18.04

MEMORY: 16G x 2 (3200MHz)

![benchmark](https://raw.githubusercontent.com/zhangyunhao116/public-data/master/skipmap-benchmark.png)

```shell
$ go test -run=NOTEST -bench=. -benchtime=100000x -benchmem -count=20 -timeout=60m  > x.txt
$ benchstat x.txt
```

```
name                                            time/op
Int64/Store/skipmap-16                           158ns ±12%
Int64/Store/sync.Map-16                          700ns ± 4%
Int64/Load50Hits/skipmap-16                     10.1ns ±14%
Int64/Load50Hits/sync.Map-16                    14.8ns ±23%
Int64/30Store70Load/skipmap-16                  50.6ns ±20%
Int64/30Store70Load/sync.Map-16                  592ns ± 7%
Int64/1Delete9Store90Load/skipmap-16            27.5ns ±13%
Int64/1Delete9Store90Load/sync.Map-16            480ns ± 4%
Int64/1Range9Delete90Store900Load/skipmap-16    34.2ns ±26%
Int64/1Range9Delete90Store900Load/sync.Map-16   1.00µs ±12%
String/Store/skipmap-16                          171ns ±15%
String/Store/sync.Map-16                         873ns ± 4%
String/Load50Hits/skipmap-16                    21.3ns ±38%
String/Load50Hits/sync.Map-16                   19.9ns ±12%
String/30Store70Load/skipmap-16                 75.6ns ±16%
String/30Store70Load/sync.Map-16                 726ns ± 5%
String/1Delete9Store90Load/skipmap-16           34.3ns ±20%
String/1Delete9Store90Load/sync.Map-16           584ns ± 5%
String/1Range9Delete90Store900Load/skipmap-16   41.0ns ±21%
String/1Range9Delete90Store900Load/sync.Map-16  1.17µs ± 8%

name                                            alloc/op
Int64/Store/skipmap-16                            112B ± 0%
Int64/Store/sync.Map-16                           128B ± 0%
Int64/Load50Hits/skipmap-16                      0.00B     
Int64/Load50Hits/sync.Map-16                     0.00B     
Int64/30Store70Load/skipmap-16                   33.0B ± 0%
Int64/30Store70Load/sync.Map-16                  81.2B ±11%
Int64/1Delete9Store90Load/skipmap-16             10.0B ± 0%
Int64/1Delete9Store90Load/sync.Map-16            57.9B ± 5%
Int64/1Range9Delete90Store900Load/skipmap-16     10.0B ± 0%
Int64/1Range9Delete90Store900Load/sync.Map-16     261B ±17%
String/Store/skipmap-16                           144B ± 0%
String/Store/sync.Map-16                          152B ± 0%
String/Load50Hits/skipmap-16                     15.0B ± 0%
String/Load50Hits/sync.Map-16                    15.0B ± 0%
String/30Store70Load/skipmap-16                  54.0B ± 0%
String/30Store70Load/sync.Map-16                 96.9B ±12%
String/1Delete9Store90Load/skipmap-16            27.0B ± 0%
String/1Delete9Store90Load/sync.Map-16           74.2B ± 4%
String/1Range9Delete90Store900Load/skipmap-16    27.0B ± 0%
String/1Range9Delete90Store900Load/sync.Map-16    257B ±10%

name                                            allocs/op
Int64/Store/skipmap-16                            3.00 ± 0%
Int64/Store/sync.Map-16                           4.00 ± 0%
Int64/Load50Hits/skipmap-16                       0.00     
Int64/Load50Hits/sync.Map-16                      0.00     
Int64/30Store70Load/skipmap-16                    0.00     
Int64/30Store70Load/sync.Map-16                   1.00 ± 0%
Int64/1Delete9Store90Load/skipmap-16              0.00     
Int64/1Delete9Store90Load/sync.Map-16             0.00     
Int64/1Range9Delete90Store900Load/skipmap-16      0.00     
Int64/1Range9Delete90Store900Load/sync.Map-16     0.00     
String/Store/skipmap-16                           4.00 ± 0%
String/Store/sync.Map-16                          5.00 ± 0%
String/Load50Hits/skipmap-16                      1.00 ± 0%
String/Load50Hits/sync.Map-16                     1.00 ± 0%
String/30Store70Load/skipmap-16                   1.00 ± 0%
String/30Store70Load/sync.Map-16                  2.00 ± 0%
String/1Delete9Store90Load/skipmap-16             1.00 ± 0%
String/1Delete9Store90Load/sync.Map-16            1.00 ± 0%
String/1Range9Delete90Store900Load/skipmap-16     1.00 ± 0%
String/1Range9Delete90Store900Load/sync.Map-16    1.00 ± 0%
```