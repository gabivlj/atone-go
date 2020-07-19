package test

import (
	"log"
	"testing"
	"time"

	"github.com/gabivlj/atone-go/atone"
)

// 2020/07/19 02:39:09 Average insertion:  136
// goos: darwin
// goarch: amd64
// pkg: github.com/gabivlj/atone-go
// BenchmarkAtone-12              1        10503182516 ns/op
// PASS
// ok      github.com/gabivlj/atone-go     10.786s
// BenchmarkAtone-12       2020/07/19 02:45:30 Average insertion:  114
// 2020/07/19 02:45:30 Average insertion:  112
// 2020/07/19 02:45:30 Average insertion:  111
// 2020/07/19 02:45:30 Average insertion:  112
// 2020/07/19 02:45:31 Average insertion:  113
// 2020/07/19 02:45:31 Average insertion:  112
// 2020/07/19 02:45:31 Average insertion:  112
// 2020/07/19 02:45:31 Average insertion:  115
// 2020/07/19 02:45:31 Average insertion:  115
// 2020/07/19 02:45:31 Average insertion:  112
// 2020/07/19 02:45:32 Average insertion:  112
// 2020/07/19 02:45:32 Average insertion:  113
// 1000000000               0.186 ns/op
// PASS
// ok      github.com/gabivlj/atone-go     2.455s
func BenchmarkAtone(b *testing.B) {
	arrMedium := int64(0)
	l := 1000000
	arr := atone.NewWithCapacity(3)
	for i := 0; i < l; i++ {
		e := time.Now()
		arr.Push(i)
		// arr.PopBack()
		arrMedium += time.Now().UnixNano() - e.UnixNano()
	}
	log.Println("Average insertion: ", arrMedium/int64(l))
	// log.Println(arr.Debug())
}

func BenchmarkThingsAtone(b *testing.B) {
	arrMedium := int64(0)
	l := 100
	arr := atone.NewWithCapacity(3)
	for i := 0; i < l; i++ {
		e := time.Now()
		arr.Push(i)
		arr.PopBack()
		arr.Push(i)
		n, _ := arr.Lookup(i)
		number, ok := n.(int)
		if !ok {
			b.Fatalf("error with number %d %s", i, arr.Debug())
			return
		}
		assert(number == i)
		arrMedium += time.Now().UnixNano() - e.UnixNano()
	}
	log.Println("Average cost in each operation: ", arrMedium/int64(l))
	// log.Println(arr.Debug())
}

// 2020/07/19 02:39:36 Average insertion:  247
// goos: darwin
// goarch: amd64
// pkg: github.com/gabivlj/atone-go
// BenchmarkStandard-12                   1        16009740195 ns/op
// PASS
// ok      github.com/gabivlj/atone-go     16.362s
// BenchmarkStandard-12            2020/07/19 02:45:34 Average insertion:  209
// 2020/07/19 02:45:35 Average insertion:  213
// 2020/07/19 02:45:35 Average insertion:  197
// 2020/07/19 02:45:35 Average insertion:  194
// 2020/07/19 02:45:36 Average insertion:  192
// 2020/07/19 02:45:36 Average insertion:  201
// 2020/07/19 02:45:36 Average insertion:  192
// 2020/07/19 02:45:36 Average insertion:  205
// 2020/07/19 02:45:37 Average insertion:  201
// 2020/07/19 02:45:37 Average insertion:  206
// 2020/07/19 02:45:37 Average insertion:  198
// 2020/07/19 02:45:38 Average insertion:  208
// 2020/07/19 02:45:38 Average insertion:  201
// 2020/07/19 02:45:38 Average insertion:  207
// 2020/07/19 02:45:38 Average insertion:  209
// 1000000000               0.282 ns/op
// PASS
// ok      github.com/gabivlj/atone-go     4.534s
func BenchmarkStandard(b *testing.B) {
	arrMedium := int64(0)
	l := 1000000
	arr2 := make([]atone.Element, 0, 3)
	for i := 0; i < l; i++ {
		e := time.Now()
		arr2 = append(arr2, i)
		arr2 = append(arr2, i)
		arr2 = append(arr2, i)
		arr2 = arr2[:len(arr2)-1]
		arrMedium += time.Now().UnixNano() - e.UnixNano()
	}
	log.Println("Average insertion: ", arrMedium/int64(l))
	// log.Println(arr2)
}

func TestPopFront(t *testing.T) {
	arr := atone.New()
	arr.Push(1)
	arr.Push(2)
	arr.Push(3)
	assert(arr.PopFront() == 1)
	assert(arr.Len() == 2)
	assert(arr.PopFront() == 2)
	assert(arr.Len() == 1)
	assert(arr.PopFront() == 3)
	assert(arr.Len() == 0)
}

func TestSwap(t *testing.T) {
	arr := atone.New()
	arr.Push(1)
	arr.Push(2)
	arr.Push(3)
	arr.Push(4)
	arr.Push(5)
	arr.Swap(0, 1)
	arr.Swap(2, 3)
	assert(arr.Get(0) == 2)
	assert(arr.Get(2) == 4)
}

func TestIter(b *testing.T) {
	nItems := 10
	arr := atone.New()
	for i := 0; i < nItems; i++ {
		arr.Push(i)
	}
	for i := 0; i < arr.Len(); i++ {
		assert(arr.Get(i) == i)
	}
	for i, el := range arr.Iter() {
		assert(el == i)
	}
	arr.ForEach(func(el atone.Element, idx int) { assert(el == idx) })
	arr.Clear()
	_, ok := arr.Lookup(0)
	assert(!ok)
}

func TestContains(b *testing.T) {
	nItems := 10
	arr := atone.New()
	for i := 0; i < nItems; i++ {
		arr.Push(i)
	}
	for i := 0; i < arr.Len(); i++ {
		assert(arr.Contains(i))
	}
	for i := 0; i < arr.Len(); i++ {
		assert(arr.ContainsCmp(i, func(el atone.Element, other atone.Element) bool { return el == other }))
	}
}

func TestModifyViaGet(b *testing.T) {
	nItems := 10
	arr := atone.New()
	for i := 0; i < nItems; i++ {
		arr.Push(i)
	}
	*arr.GetRef(0) = 2 + 3
	assert(arr.Get(0) == 5)
}

func assert(cond bool) {
	if !cond {
		panic("condition not met")
	}
}
