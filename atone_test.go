package test

import (
	"log"
	"testing"
	"time"

	"github.com/gabivlj/atone-go/atone"
)

func BenchmarkThingsAtone(b *testing.B) {
	arrMedium := int64(0)
	l := 10000000
	arr := atone.NewWithCapacity[int](3)
	for i := 0; i < l; i++ {
		e := time.Now()
		arr.Push(i)
		arr.PopBack()
		arr.Push(i)
		n, _ := arr.Lookup(i)
		number := n
		assert(number == i)
		arrMedium += time.Now().UnixNano() - e.UnixNano()
	}
	log.Println("Average cost in each operation: ", arrMedium/int64(l))
	// log.Println(arr.Debug())
}

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
	l := 10000000
	m := int64(0)
	arr := atone.New[int]()
	stats := make([]int64, 0, l)
	for i := 0; i < l; i++ {
		e := time.Now()
		arr.Push(i)
		// arr.PopBack()
		t := time.Now().UnixNano() - e.UnixNano()
		arrMedium += t
		m = max(t, m)
		if t > 300 {
			stats = append(stats, t)
		}
	}
	log.Println("Average insertion: ", arrMedium/int64(l))
	log.Println("Max insertion: ", int64(m))
	log.Println("Stats (Number of times an insert took a lot of time): ", len(stats))
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
	l := 10000000
	m := int64(0)
	arr2 := make([]int, 0)
	stats := make([]int64, 0, l)
	for i := 0; i < l; i++ {
		e := time.Now()
		arr2 = append(arr2, i)
		t := time.Now().UnixNano() - e.UnixNano()
		arrMedium += t
		m = max(t, m)
		if t > 200 {
			stats = append(stats, t)
		}
	}
	log.Println("Average insertion: ", arrMedium/int64(l))
	log.Println("Max insertion: ", int64(m))
	log.Println("Stats (Number of times an insert took a lot of time): ", len(stats))
	// log.Println(arr2)
}

// 196793000
// 88722000

func TestPopFront(t *testing.T) {
	arr := atone.New[int]()
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
	arr := atone.New[int]()
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
	arr := atone.New[int]()
	for i := 0; i < nItems; i++ {
		arr.Push(i)
	}
	for i := 0; i < arr.Len(); i++ {
		assert(arr.Get(i) == i)
	}
	for i, el := range arr.Iter() {
		assert(el == i)
	}
	arr.ForEach(func(el int, idx int) { assert(el == idx) })
	arr.Clear()
	_, ok := arr.Lookup(0)
	assert(!ok)
}

func TestContains(b *testing.T) {
	nItems := 10
	arr := atone.New[int]()
	for i := 0; i < nItems; i++ {
		arr.Push(i)
	}
	for i := 0; i < arr.Len(); i++ {
		assert(arr.Contains(i, func(el int) bool { return i == el }))
	}
	for i := 0; i < arr.Len(); i++ {
		assert(arr.ContainsCmp(i, func(el int, other int) bool { return el == other }))
	}
}

func TestModifyViaGet(b *testing.T) {
	nItems := 10
	arr := atone.New[int]()
	for i := 0; i < nItems; i++ {
		arr.Push(i)
	}
	*arr.GetRef(0) = 2 + 3
	assert(arr.Get(0) == 5)
}

func BenchmarkFindMulti(b *testing.B) {
	nItems := 982771
	arr := atone.New[int]()
	for i := 0; i < nItems; i++ {
		arr.Push(i)
	}
	assert(arr.FindMultithreaded(355523, func(element int) bool { return element == 355523 }) == 355523)
}

func BenchmarkFind(b *testing.B) {
	nItems := 982771
	arr := atone.New[int]()
	for i := 0; i < nItems; i++ {
		arr.Push(i)
	}
	assert(arr.Find(355523, func(element int) bool { return element == 355523 }) == 355523)
}

func BenchmarkFind02(b *testing.B) {
	nItems := 1000000
	arr := atone.New[int]()
	for i := 0; i < nItems; i++ {
		arr.Push(i)
	}
	assert(arr.Find02(355523, func(element int) bool { return element == 355523 }) == 355523)
}

func BenchmarkInsertAtone(b *testing.B) {
	arrMedium := int64(0)
	nItems := 10000
	arr := atone.New[int]()
	for i := 0; i < nItems; i++ {
		e := time.Now()
		arr.Insert(i)
		assert(arr.Get(0) == i)
		arrMedium += time.Now().UnixNano() - e.UnixNano()
	}
	log.Println("Average insertion: ", arrMedium/int64(nItems))
	for i := 0; i < nItems; i++ {
		assert(arr.Get(i) == nItems-1-i)
	}
}

func BenchmarkInsert(b *testing.B) {
	arrMedium := int64(0)
	nItems := 10000
	arr := make([]int, 0)
	for i := 0; i < nItems; i++ {
		e := time.Now()
		arr = append([]int{1}, arr...)
		assert(arr[0] == i)
		arrMedium += time.Now().UnixNano() - e.UnixNano()
	}
	log.Println("Average insertion: ", arrMedium/int64(nItems))
	for i := 0; i < nItems; i++ {
		assert(arr[i] == nItems-1-i)
	}
}

func TestReverse(t *testing.T) {
	nItems := 26
	arr := atone.New[int]()
	for i := 0; i < nItems; i++ {
		arr.Push(i)
	}
	arr.Reverse()
	for i := 0; i < nItems; i++ {
		assert(arr.Get(i) == nItems-1-i)
	}
}

func TestReserve(t *testing.T) {
	nItems := 17
	arr := atone.NewWithCapacity[int](20)
	for i := 0; i < nItems; i++ {
		arr.Push(i)
	}
	// assert(arr.Capacity() >= 25)
}

func TestTruncate(t *testing.T) {
	nItems := 17
	arr := atone.NewWithCapacity[int](0)
	for i := 0; i < nItems; i++ {
		arr.Push(i)
	}

	for i := 0; i < nItems; i++ {
		arr.Truncate(nItems - i)
		assert(arr.Len() == nItems-i)
	}
}

func TestSlice(t *testing.T) {
	nItems := 17
	arr := atone.NewWithCapacity[int](0)
	for i := 0; i < nItems; i++ {
		arr.Push(i)
	}
	for i := 0; i < 14; i++ {
		elements := arr.Slice(i, 13)
		assert(len(elements) == 14-i-1)
		for j := range elements {
			assert(j+i == elements[j])
		}
	}
	assert(len(arr.Array()) == nItems)
	for i, el := range atone.From(arr.Array()).Iter() {
		assert(arr.Get(i) == el)
	}
}

func assert(cond bool) {
	if !cond {
		panic("condition not met")
	}
}

func max(n, n2 int64) int64 {
	if n > n2 {
		return n
	}
	return n2
}
