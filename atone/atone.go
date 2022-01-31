///
/// This is an implementation of atone in Golang
/// originally made by @jonhoo in Rust. Original repository: https://github.com/jonhoo/atone
///
/// Implementation by @gabivlj. Free to use and contribute by anyone.
///
/// This is currently under development and it's not yet very optimized
///
/// atone.go has GC overhead.

package atone

import (
	"fmt"
	"log"
)

// Debug is true if we should print debug statements
var Debug = false

// Vec is the implementation of an atone vector just like in https://github.com/jonhoo/atone, there you will find what is so special about this implementation
type Vec[T any] struct {
	oldHead []T
	newTail []T
}

// NItemsToMoveOnEachInsert is the number of items we move on each insert, between 4-8 the performance doesn't have much difference
const NItemsToMoveOnEachInsert = 5

// New returns a new atone Vec
func New[T any]() *Vec[T] {
	return &Vec[T]{
		newTail: make([]T, 0, 0),
		oldHead: nil,
	}
}

// From returns a new Vec from a Slice
func From[T any](elements []T) *Vec[T] {
	return &Vec[T]{
		newTail: elements,
		oldHead: nil,
	}
}

// NewWithCapacity is the equivalent of doing make([]T, 0, capacity)
func NewWithCapacity[T any](capacity uint64) *Vec[T] {
	return &Vec[T]{
		newTail: make([]T, 0, capacity),
		oldHead: nil,
	}
}

// returns the oldHead len
func (v *Vec[T]) oldLen() int {
	if v.oldHead == nil {
		return 0
	}
	return len(v.oldHead)
}

// Lookup returns an element, the boolean is false if the element does not exist.
func (v *Vec[T]) Lookup(index int) (T, bool) {
	var defaul T
	if index < v.oldLen() {
		return v.oldHead[index], true
	}
	offset := index - v.oldLen()
	if offset >= len(v.newTail) || offset < 0 {
		return defaul, false
	}
	return v.newTail[offset], true
}

// Get returns the element in the specified index, can panic if it is outofbounds, if you don't want to panic on get, use Lookup
func (v *Vec[T]) Get(index int) T {
	if index < v.oldLen() {
		return v.oldHead[index]
	}

	offset := index - v.oldLen()
	return v.newTail[offset]
}

// GetRef returns FOR SURE a pointer to the element even though it is a stack element like int
func (v *Vec[T]) GetRef(index int) *T {
	if index < v.oldLen() {
		return &v.oldHead[index]
	}

	offset := index - v.oldLen()
	return &v.newTail[offset]
}

// Find02 tries to find not doing a continuous loop
func (v *Vec[T]) Find02(el T, cb func(element T) bool) int {
	bigger := v.newTail
	smaller := v.oldHead
	if len(v.newTail) < v.oldLen() {
		bigger, smaller = smaller, bigger
	}
	for i := range bigger {
		if cb(bigger[i]) {
			return i
		}
		if i < len(smaller) && cb(smaller[i]) {
			return i
		}
	}
	return -1
}

// Find finds doing a lookup in head and then in tail
func (v *Vec[T]) Find(el T, cb func(element T) bool) int {
	if v.oldHead != nil {
		for i := range v.oldHead {
			if cb(v.oldHead[i]) {
				return i
			}
		}
	}

	for i := range v.newTail {
		if cb(v.newTail[i]) {
			return i
		}
	}

	return -1
}

func find[T any](els []T, el T, cb func(T) bool) int {
	for i := range els {
		if cb(els[i]) {
			return i
		}
	}
	return -1
}

// FindMultithreaded finds an element with multithreading (with a lot of elements 1000000+)
func (v *Vec[T]) FindMultithreaded(el T, cb func(T) bool) int {
	channel := make(chan int, 2)
	go func() {
		channel <- find(v.oldHead, el, cb)
	}()
	go func() {
		channel <- find(v.newTail, el, cb)
	}()

	target := 0
	for found := range channel {
		target++
		if found != -1 {
			return found
		}
		if target >= 2 {
			break
		}
	}

	close(channel)

	return -1
}

// Insert .
func (v *Vec[T]) Insert(el T) {

	if len(v.newTail) == cap(v.newTail) {
		v.grow(1)
		v.Insert(el)
		return
	}
	if v.oldLen() == 0 {
		els := make([]T, 0, cap(v.newTail)+1)
		els = append(els, el)
		v.newTail = append(els, v.newTail...)
		return
	}
	// storage for sufficient elements in the new tail. maybe with a better implementaiton we could jump this
	els := make([]T, 0, cap(v.newTail)+1)
	// insert the popped element
	els[0] = v.oldHead[len(v.oldHead)-1]
	// copy rest
	v.newTail = append(els, v.newTail...)
	// length
	oldLen := len(v.oldHead)
	// append the element to the old head
	v.oldHead = append(v.oldHead, el)
	// append rest
	v.oldHead = append(v.oldHead, v.oldHead...)
	// jump
	v.oldHead = v.oldHead[oldLen:]
	if v.oldLen() != 0 {
		v.carry()
	}
}

// Swap swaps elements in the structure
func (v *Vec[T]) Swap(i int, j int) {
	iIsInOldHead := i < v.oldLen()
	jIsInOldHead := j < v.oldLen()

	if iIsInOldHead == jIsInOldHead {
		if iIsInOldHead {
			v.oldHead[i], v.oldHead[j] = v.oldHead[j], v.oldHead[i]
			return
		}
		l := v.oldLen()
		v.newTail[i-l], v.newTail[j-l] = v.newTail[j-l], v.newTail[i-l]
		return
	}

	if !iIsInOldHead {
		l := v.oldLen()
		v.oldHead[i], v.newTail[j-l] = v.newTail[j-l], v.oldHead[i]
		return
	}

	v.oldHead[j], v.newTail[i] = v.newTail[i], v.oldHead[j]
}

func reverseSlice[T any](s []T) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// Reverse inplace the array and empties the old head
func (v *Vec[T]) Reverse() {
	reverseSlice(v.newTail)
	if v.oldHead != nil {
		for i := range v.oldHead {
			v.newTail = append(v.newTail, v.oldHead[v.oldLen()-i-1])
		}
		v.oldHead = nil
	}
}

// Reserve the desired size inmemory to let space for nElements, it might reserve more memory than necessary for leaving space for more items for carry()
func (v *Vec[T]) Reserve(nElements int) {
	if v.oldLen() > 0 {
		v.carryAll()
	}

	v.grow(nElements)
}

// Capacity is the equivalent of cap(elements)
func (v *Vec[T]) Capacity() int {
	return cap(v.newTail)
}

// Shrink ; the capacity will remain to atleast the length of the array (TODO)
func (v *Vec[T]) Shrink(minCapacity int) {
	// Calculate the min. number of elements we need to reserve
	need := len(v.newTail)
	if v.oldHead != nil {
		oldLen := v.oldLen()
		need += oldLen
		need += (oldLen + NItemsToMoveOnEachInsert - 1) / NItemsToMoveOnEachInsert
	} else if minCapacity <= need {
		panic("todo")
	}
	panic("todo")
}

// Truncate only will mantain only the first 'n' elements in the array and the rest will be free'd
func (v *Vec[T]) Truncate(n int) {
	if n <= v.oldLen() {
		v.newTail = append(v.newTail[:0], v.oldHead[:n]...)
		if n == v.oldLen() {
			v.oldHead = nil
		} else {
			v.oldHead = v.oldHead[n:]
		}
		return
	}
	maintain := n - v.oldLen()
	v.newTail = append(v.oldHead, v.newTail[:maintain]...)
}

// Len returns the number of elements stored in the array
func (v *Vec[T]) Len() int {
	return v.oldLen() + len(v.newTail)
}

// IsEmpty returns if there is any element in the array or not
func (v *Vec[T]) IsEmpty() bool {
	return v.Len() == 0
}

// Clear empties the array
func (v *Vec[T]) Clear() {
	v.oldHead = nil
	v.newTail = v.newTail[:0]
}

// Contains returns true if the element is inside the array
func (v *Vec[T]) Contains(el T, cb func(T) bool) bool {
	return v.Find(el, cb) > -1
}

// ContainsCmp returns true if the element is inside the array, will use the cmp func
func (v *Vec[T]) ContainsCmp(el T, cmp func(arrayElement T, el T) bool) bool {
	bigger := v.newTail
	smaller := v.oldHead
	if len(v.newTail) < v.oldLen() {
		bigger, smaller = smaller, bigger
	}
	for i := range bigger {
		if cmp(bigger[i], el) {
			return true
		}
		if smaller != nil && i < len(smaller) && cmp(smaller[i], el) {
			return true
		}
	}
	return false
}

// First returns the first element of the array, returns null if it is empty
func (v *Vec[T]) First() T {
	var t T
	if v.oldLen() > 0 {
		return v.oldHead[0]
	}
	if len(v.newTail) > 0 {
		return v.newTail[0]
	}
	return t
}

// Last returns the last element of the array, returns null if it is empty
func (v *Vec[T]) Last() T {
	var t T
	if len(v.newTail) > 0 {
		return v.newTail[len(v.newTail)-1]
	}
	oldLen := v.oldLen()
	if oldLen > 0 {
		return v.oldHead[oldLen]
	}
	return t
}

// PopFront pops the first element of the array, returns null if the array is empty
func (v *Vec[T]) PopFront() T {
	var t T
	if v.oldLen() > 0 {
		popped := v.oldHead[0]
		v.oldHead = v.oldHead[1:]
		return popped
	}
	if len(v.newTail) > 0 {
		popped := v.newTail[0]
		v.newTail = v.newTail[1:]
		return popped
	}
	return t
}

// PopBack pops the last element of the array, returns null if the array is empty
func (v *Vec[T]) PopBack() T {
	var t T
	if len(v.newTail) > 0 {
		popped := v.newTail[len(v.newTail)-1]
		v.newTail = v.newTail[:len(v.newTail)-1]
		return popped
	}
	oldL := v.oldLen()
	if oldL > 0 {
		popped := v.oldHead[oldL-1]
		v.oldHead = v.oldHead[:oldL-1]
		return popped
	}
	return t
}

// Pop same as PopBack
func (v *Vec[T]) Pop() T {
	return v.PopBack()
}

// Iter generates an array of elements (allocates space for the iteration)
func (v *Vec[T]) Iter() []T {
	elements := make([]T, 0, v.Len())
	if v.oldHead != nil {
		elements = append(elements, v.oldHead...)
	}
	elements = append(elements, v.newTail...)
	return elements
}

// Slice generates a slice slicing the array from start to end (end is not inclusive and start is)
func (v *Vec[T]) Slice(start, end int) []T {
	elements := make([]T, 0, end-start)
	if v.oldLen() > start {
		newEnd := end
		if newEnd > v.oldLen() {
			newEnd = v.oldLen()
		}
		elements = append(elements, v.oldHead[start:newEnd]...)
		if newEnd > v.oldLen() {
			elements = append(elements, v.newTail[:newEnd-v.oldLen()]...)
		}
		return elements
	}
	elements = append(elements, v.newTail[start:end]...)
	return elements
}

// Array creates a slice of this array
func (v *Vec[T]) Array() []T {
	return v.Slice(0, v.Len())
}

// SliceThis returns a new Vec with the specified slice (end non inclusive and start is inclusive)
func (v *Vec[T]) SliceThis(start, end int) *Vec[T] {
	return From(v.Slice(start, end))
}

// ForEach iterates through the array doing a callback to the passed function
func (v *Vec[T]) ForEach(fn func(el T, index int)) {
	if v.oldLen() > 0 {
		for i := range v.oldHead {
			fn(v.oldHead[i], i)
		}
	}
	for i := range v.newTail {
		fn(v.newTail[i], i)
	}
}

// Push pushes back an element into the array
func (v *Vec[T]) Push(el T) {
	if cap(v.newTail) == len(v.newTail) {
		v.grow(1)
		v.Push(el)
		return
	}

	v.newTail = append(v.newTail, el)
	if v.oldLen() != 0 {
		v.carry()
	}
}

// Append is the equivalent of doing append(elements, toAppend...)
func (v *Vec[T]) Append(el ...T) {
	for _, e := range el {
		v.Push(e)
	}
}

func (v *Vec[T]) carry() {
	if v.oldLen() == 0 {
		v.oldHead = nil
		return
	}
	lenOld := v.oldLen()
	calc := max(lenOld-NItemsToMoveOnEachInsert, 0)
	v.newTail = append(v.oldHead[calc:], v.newTail...)
	v.oldHead = v.oldHead[:calc]
	if v.oldLen() == 0 {
		v.oldHead = nil
		return
	}

}

func (v *Vec[T]) carryAll() {
	if v.oldLen() == 0 {
		v.oldHead = nil
		return
	}
	v.newTail = append(v.oldHead[0:], v.newTail...)
	v.oldHead = nil
}

const pushMultiplierOldVector = 2

func (v *Vec[T]) grow(growFactor int) {
	// assertDebug(v.oldLen() == 0)
	// Original repo comments
	// We need to grow the Vec by at least a factor of (R + 1)/R to ensure that
	// the new Vec won't _also_ grow while we're still moving items from the old
	// one.
	//
	// Here's how we get to len * (R + 1)/R:
	//  - We need to move another len items
	need := len(v.newTail)
	//  - We move R items on each push, so to move len items takes
	//    len / R pushes (rounded up!)
	//  - Since we want to round up, we pull the old +R-1 trick
	pushes := (need + NItemsToMoveOnEachInsert - 1) / NItemsToMoveOnEachInsert
	//  - That's len + len/R
	//    Which is == R*len/R + len/R
	//    Which is == ((R+1)*len)/R
	//    Which is == len * (R+1)/R
	//  - We don't actually use that formula because of integer division.
	// We also need to make sure we can fit the additional capacity required for `extra`.
	// Normally, that'll be handled by `pushes`, but not always!
	add := max(pushes, growFactor)
	elements := make([]T, 0, cap(v.newTail)+add+pushes+need)
	v.oldHead = make([]T, 0, add+pushes+need*pushMultiplierOldVector)
	v.oldHead = append(v.oldHead, v.newTail...)

	v.newTail = elements
}

func max(n, n2 int) int {
	if n > n2 {
		return n
	}
	return n2
}

func assertDebug(cond bool) {
	if !Debug || cond {
		return
	}
	log.Println("the condition is false!")
}

// Debug the vec
func (v *Vec[T]) Debug() string {
	return fmt.Sprintf("Old: %v \n", v.oldHead) + fmt.Sprintf("New: %v", v.newTail)
}
