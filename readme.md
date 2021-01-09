<p align="center">
  <img src="https://raw.githubusercontent.com/ZYunH/public-data/master/skipmap-logo.png"/>
</p>

## Introduction

skipmap is a high-performance concurrent map based on skip list. In typical pattern(one million operations, 90%LOAD 9%STORE 1%DELETE), the skipmap up to 3x ~ 10x faster than the built-in sync.Map.

The main idea behind the skipmap is [A Simple Optimistic Skiplist Algorithm](<https://people.csail.mit.edu/shanir/publications/LazySkipList.pdf>).

Different from the sync.Map, the items in the skipmap are always sorted, and the `Load` and `Range` operations are wait-free (A goroutine is guaranteed to complete a operation as long as it keeps taking steps, regardless of the activity of other goroutines).



## Features

- Concurrent safe API with high-performance.
- Wait-free Load and Range operations.
- Sorted items.



## When should you use skipmap

In these situations, `skipmap` is better

- **Sorted elements is needed**.
- **Concurrent calls multiple operations**. such as use both `Load` and `Store` at the same time.

In these situations, `sync.Map` is better

- Only one goroutine use the map for most of the time, such as insert a batch of elements and then use only `Load` (use built-in map is even better).



## QuickStart

See [Go doc](https://pkg.go.dev/github.com/ZYunH/skipmap) for more information.

```go
package main

import (
	"fmt"

	"github.com/zhangyunhao116/skipmap"
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
Store/skipmap-16                                267ns ± 5%
Store/sync.Map-16                               675ns ± 6%
Load100Hits/skipmap-16                         15.2ns ± 6%
Load100Hits/sync.Map-16                        16.0ns ±11%
Load50Hits/skipmap-16                          15.6ns ± 7%
Load50Hits/sync.Map-16                         14.2ns ±18%
LoadNoHits/skipmap-16                          16.7ns ±21%
LoadNoHits/sync.Map-16                         13.1ns ±18%
50Store50Load/skipmap-16                        151ns ±38%
50Store50Load/sync.Map-16                       568ns ± 2%
30Store70Load/skipmap-16                       95.2ns ±43%
30Store70Load/sync.Map-16                       584ns ± 4%
1Delete9Store90Load/skipmap-16                 46.0ns ±11%
1Delete9Store90Load/sync.Map-16                 505ns ± 4%
1Range9Delete90Store900Load/skipmap-16         52.5ns ± 8%
1Range9Delete90Store900Load/sync.Map-16        1.15µs ±18%
StringStore/skipmap-16                          321ns ± 7%
StringStore/sync.Map-16                         872ns ± 4%
StringLoad50Hits/skipmap-16                    28.6ns ± 6%
StringLoad50Hits/sync.Map-16                   18.7ns ± 8%
String30Store70Load/skipmap-16                  125ns ± 5%
String30Store70Load/sync.Map-16                 746ns ± 6%
String1Delete9Store90Load/skipmap-16           56.9ns ± 8%
String1Delete9Store90Load/sync.Map-16           619ns ± 3%
String1Range9Delete90Store900Load/skipmap-16   64.8ns ±24%
String1Range9Delete90Store900Load/sync.Map-16  1.46µs ±20%

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
50Store50Load/sync.Map-16                       65.2B ± 1%
30Store70Load/skipmap-16                        32.0B ± 0%
30Store70Load/sync.Map-16                       74.4B ± 3%
1Delete9Store90Load/skipmap-16                  9.00B ± 0%
1Delete9Store90Load/sync.Map-16                 55.4B ± 3%
1Range9Delete90Store900Load/skipmap-16          9.00B ± 0%
1Range9Delete90Store900Load/sync.Map-16          286B ± 9%
StringStore/skipmap-16                           138B ± 0%
StringStore/sync.Map-16                          152B ± 0%
StringLoad50Hits/skipmap-16                     3.00B ± 0%
StringLoad50Hits/sync.Map-16                    3.00B ± 0%
String30Store70Load/skipmap-16                  52.0B ± 0%
String30Store70Load/sync.Map-16                 96.6B ±13%
String1Delete9Store90Load/skipmap-16            26.0B ± 0%
String1Delete9Store90Load/sync.Map-16           72.3B ± 1%
String1Range9Delete90Store900Load/skipmap-16    26.0B ± 0%
String1Range9Delete90Store900Load/sync.Map-16    333B ±23%
```