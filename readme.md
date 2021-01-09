<p align="center">
  <img src="https://raw.githubusercontent.com/ZYunH/public-data/master/skipmap-logo.png"/>
</p>

## Introduction

skipmap is a high-performance concurrent map based on skip list. In typical pattern(one million operations, 90%CONTAINS 9%INSERT 1%DELETE), the skipmap up to 3x ~ 10x faster than the built-in sync.Map.

The main idea behind the skipmap is [A Simple Optimistic Skiplist Algorithm](<https://people.csail.mit.edu/shanir/publications/LazySkipList.pdf>).

Different from the sync.Map, the items in the skipmap are always sorted, and the `Load` and `Range` operations are wait-free (A goroutine is guaranteed to complete a operation as long as it keeps taking steps, regardless of the activity of other goroutines).



## Features

- Concurrent safe API with high-performance.
- Wait-free Contains and Range operations.
- Sorted items.



## When should you use skipmap

In these situations, `skipmap` is better

- **Sorted elements is needed**.
- **Concurrent calls multiple operations**. such as use both `Load` and `Store` at the same time.

In these situations, `sync.Map` is better

- Only one goroutine use the map for most of the time, such as insert a batch of elements and then use only `Load` (use built-in map is even better), 



## QuickStart

See [Go doc](https://pkg.go.dev/github.com/ZYunH/skipmap) for more information.

```go
package main

import (
	"fmt"

	"github.com/ZYunH/skipmap"
)

func main() {
	l := skipmap.NewInt()

	for _, v := range []int{10, 12, 15} {
		l.Store(v, v+100)
	}

	v, ok := l.Load(10)
	if ok {
		fmt.Println("skipmap load 10 with value ", v)
	}

	l.Range(func(key int, value interface{}) bool {
		fmt.Println("skipmap range found ", key, value)
		return true
	})

	l.Delete(15)
	fmt.Printf("skipmap contains %d items\r\n", l.Len())
}

```



## Benchmark

Go version: go1.15.6 linux/amd64

CPU: AMD 3700x(8C16T), running at 3.6GHz

OS: ubuntu 18.04

MEMORY: 16G x 2 (3200MHz)

![benchmark](https://raw.githubusercontent.com/ZYunH/public-data/master/skipmap-benchmark.png)

```shell
$ go test -run=NOTEST -bench=. -benchtime=100000x -benchmem -count=10 -timeout=60m  > x.txt
$ benchstat x.txt
```

```
name                                           time/op
Store/skipmap-16                                287ns ±21%
Store/sync.Map-16                               684ns ± 5%
Load100Hits/skipmap-16                         15.2ns ±14%
Load100Hits/sync.Map-16                        15.9ns ±18%
Load50Hits/skipmap-16                          15.5ns ± 4%
Load50Hits/sync.Map-16                         14.5ns ±16%
LoadNoHits/skipmap-16                          17.2ns ±22%
LoadNoHits/sync.Map-16                         12.2ns ±11%
50Store50Load/skipmap-16                        149ns ±13%
50Store50Load/sync.Map-16                       569ns ± 7%
30Store70Load/skipmap-16                       86.5ns ± 9%
30Store70Load/sync.Map-16                       601ns ± 6%
1Delete9Store90Load/skipmap-16                 47.5ns ±14%
1Delete9Store90Load/sync.Map-16                 509ns ± 3%
1Range9Delete90Store900Load/skipmap-16         55.8ns ±11%
1Range9Delete90Store900Load/sync.Map-16        1.13µs ±11%
StringStore/skipmap-16                          368ns ±10%
StringStore/sync.Map-16                         881ns ± 4%
StringLoad50Hits/skipmap-16                    29.8ns ±14%
StringLoad50Hits/sync.Map-16                   20.3ns ±26%
String30Store70Load/skipmap-16                  130ns ± 7%
String30Store70Load/sync.Map-16                 751ns ± 5%
String1Delete9Store90Load/skipmap-16           40.4ns ±22%
String1Delete9Store90Load/sync.Map-16           446ns ± 7%
String1Range9Delete90Store900Load/skipmap-16   68.5ns ±21%
String1Range9Delete90Store900Load/sync.Map-16  1.32µs ±13%

name                                           alloc/op
Store/skipmap-16                                 107B ± 0%
Store/sync.Map-16                                128B ± 0%
Load100Hits/skipmap-16                          0.00B     
Load100Hits/sync.Map-16                         0.00B     
Load50Hits/skipmap-16                           0.00B     
Load50Hits/sync.Map-16                          0.00B     
LoadNoHits/skipmap-16                           0.00B     
LoadNoHits/sync.Map-16                          0.00B     
50Store50Load/skipmap-16                        53.0B ± 0%
50Store50Load/sync.Map-16                       66.0B ± 5%
30Store70Load/skipmap-16                        32.0B ± 0%
30Store70Load/sync.Map-16                       82.7B ±12%
1Delete9Store90Load/skipmap-16                  9.00B ± 0%
1Delete9Store90Load/sync.Map-16                 55.4B ± 3%
1Range9Delete90Store900Load/skipmap-16          9.00B ± 0%
1Range9Delete90Store900Load/sync.Map-16          290B ±12%
StringStore/skipmap-16                           138B ± 0%
StringStore/sync.Map-16                          152B ± 0%
StringLoad50Hits/skipmap-16                     3.00B ± 0%
StringLoad50Hits/sync.Map-16                    3.00B ± 0%
String30Store70Load/skipmap-16                  52.0B ± 0%
String30Store70Load/sync.Map-16                 96.6B ±12%
String1Delete9Store90Load/skipmap-16            16.6B ± 4%
String1Delete9Store90Load/sync.Map-16           57.7B ± 5%
String1Range9Delete90Store900Load/skipmap-16    26.0B ± 0%
String1Range9Delete90Store900Load/sync.Map-16    298B ±21%
```