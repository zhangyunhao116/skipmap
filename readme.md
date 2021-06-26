<p align="center">
  <img src="https://raw.githubusercontent.com/zhangyunhao116/public-data/master/skipmap-logo.png"/>
</p>

## Introduction

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
Int64/Store/skipmap-16                           187ns ±11%
Int64/Store/sync.Map-16                          708ns ± 5%
Int64/Load50Hits/skipmap-16                     13.7ns ±22%
Int64/Load50Hits/sync.Map-16                    14.4ns ±26%
Int64/30Store70Load/skipmap-16                  68.8ns ±14%
Int64/30Store70Load/sync.Map-16                  585ns ± 6%
Int64/1Delete9Store90Load/skipmap-16            37.7ns ±30%
Int64/1Delete9Store90Load/sync.Map-16            479ns ± 3%
Int64/1Range9Delete90Store900Load/skipmap-16    43.2ns ±21%
Int64/1Range9Delete90Store900Load/sync.Map-16    974ns ±19%
String/Store/skipmap-16                          207ns ±13%
String/Store/sync.Map-16                         879ns ± 4%
String/Load50Hits/skipmap-16                    21.1ns ±21%
String/Load50Hits/sync.Map-16                   21.1ns ±39%
String/30Store70Load/skipmap-16                 86.6ns ± 8%
String/30Store70Load/sync.Map-16                 722ns ± 7%
String/1Delete9Store90Load/skipmap-16           43.2ns ±17%
String/1Delete9Store90Load/sync.Map-16           577ns ± 3%
String/1Range9Delete90Store900Load/skipmap-16   47.3ns ± 9%
String/1Range9Delete90Store900Load/sync.Map-16  1.22µs ±16%

name                                            alloc/op
Int64/Store/skipmap-16                            106B ± 0%
Int64/Store/sync.Map-16                           128B ± 0%
Int64/Load50Hits/skipmap-16                      0.00B     
Int64/Load50Hits/sync.Map-16                     0.00B     
Int64/30Store70Load/skipmap-16                   31.3B ± 2%
Int64/30Store70Load/sync.Map-16                  80.8B ±11%
Int64/1Delete9Store90Load/skipmap-16             9.00B ± 0%
Int64/1Delete9Store90Load/sync.Map-16            58.2B ± 5%
Int64/1Range9Delete90Store900Load/skipmap-16     9.00B ± 0%
Int64/1Range9Delete90Store900Load/sync.Map-16     245B ±24%
String/Store/skipmap-16                           138B ± 0%
String/Store/sync.Map-16                          152B ± 0%
String/Load50Hits/skipmap-16                     15.0B ± 0%
String/Load50Hits/sync.Map-16                    15.0B ± 0%
String/30Store70Load/skipmap-16                  52.0B ± 0%
String/30Store70Load/sync.Map-16                 97.3B ±12%
String/1Delete9Store90Load/skipmap-16            26.0B ± 0%
String/1Delete9Store90Load/sync.Map-16           74.4B ± 2%
String/1Range9Delete90Store900Load/skipmap-16    26.0B ± 0%
String/1Range9Delete90Store900Load/sync.Map-16    275B ±18%

name                                            allocs/op
Int64/Store/skipmap-16                            4.00 ± 0%
Int64/Store/sync.Map-16                           4.00 ± 0%
Int64/Load50Hits/skipmap-16                       0.00     
Int64/Load50Hits/sync.Map-16                      0.00     
Int64/30Store70Load/skipmap-16                    1.00 ± 0%
Int64/30Store70Load/sync.Map-16                   1.00 ± 0%
Int64/1Delete9Store90Load/skipmap-16              0.00     
Int64/1Delete9Store90Load/sync.Map-16             0.00     
Int64/1Range9Delete90Store900Load/skipmap-16      0.00     
Int64/1Range9Delete90Store900Load/sync.Map-16     0.00     
String/Store/skipmap-16                           5.00 ± 0%
String/Store/sync.Map-16                          5.00 ± 0%
String/Load50Hits/skipmap-16                      1.00 ± 0%
String/Load50Hits/sync.Map-16                     1.00 ± 0%
String/30Store70Load/skipmap-16                   2.00 ± 0%
String/30Store70Load/sync.Map-16                  2.00 ± 0%
String/1Delete9Store90Load/skipmap-16             1.00 ± 0%
String/1Delete9Store90Load/sync.Map-16            1.00 ± 0%
String/1Range9Delete90Store900Load/skipmap-16     1.00 ± 0%
String/1Range9Delete90Store900Load/sync.Map-16    1.00 ± 0%
```