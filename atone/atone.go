///
/// This is an implementation of atone in Golang
/// originally made by @jonhoo in Rust. Original repository: https://github.com/jonhoo/atone
///
/// Implementation by @gabivlj. Free to use and contribute by anyone.
///
/// This is currently under development and it's not yet very optimized

package atone

import (
	"fmt"
	"log"
)

// Debug is true if we should print debug statements
var Debug = false

// Element interface
type Element interface{}

// Vec is the implementation of an atone vector just like in https://github.com/jonhoo/atone, there you will find what is so special about this implementation
type Vec struct {
	oldHead []Element
	newTail []Element
}

// NItemsToMoveOnEachInsert is the number of items we move on each insert, recommended values: 1, 4, 6....
const NItemsToMoveOnEachInsert = 6

// New returns a new atone Vec
func New() Vec {
	return Vec{
		newTail: make([]Element, 0, 0),
		oldHead: nil,
	}
}

// NewWithCapacity is the equivalent of doing make([]Element, 0, capacity)
func NewWithCapacity(capacity uint64) Vec {
	return Vec{
		newTail: make([]Element, 0, capacity),
		oldHead: nil,
	}
}

// returns the oldHead len
func (v *Vec) oldLen() int {
	if v.oldHead == nil {
		return 0
	}
	return len(v.oldHead)
}

// Get returns an element, the boolean is false if the element does not exist.
func (v *Vec) Get(index int) (Element, bool) {
	if index < v.oldLen() {
		return v.oldHead[index], true
	}
	offset := index - v.oldLen()
	if offset >= len(v.newTail) || offset < 0 {
		return nil, false
	}
	return v.newTail[offset], true
}

// Swap swaps elements in the structure
func (v *Vec) Swap(i int, j int) {
	iIsInOldHead := i < v.oldLen()
	jIsInOldHead := j < v.oldLen()

	if iIsInOldHead == jIsInOldHead {
		if iIsInOldHead {
			v.oldHead[i], v.oldHead[j] = v.oldHead[j], v.oldHead[i]
			return
		}
		v.newTail[i], v.newTail[j] = v.newTail[j], v.newTail[i]
		return
	}

	if !iIsInOldHead {
		v.oldHead[i], v.newTail[j] = v.newTail[j], v.oldHead[i]
		return
	}

	v.oldHead[j], v.newTail[i] = v.newTail[i], v.oldHead[j]
}

// func Reverse todo

// Capacity is the equivalent of cap(elements)
func (v *Vec) Capacity() int {
	return cap(v.newTail)
}

// Shrink ; the capacity will remain to atleast the length of the array (TODO)
func (v *Vec) Shrink(minCapacity int) {
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
func (v *Vec) Truncate(n int) {}

// Len returns the number of elements stored in the array
func (v *Vec) Len() int {
	return v.oldLen() + len(v.newTail)
}

// IsEmpty returns if there is any element in the array or not
func (v *Vec) IsEmpty() bool {
	return v.Len() == 0
}

// Clear empties the array
func (v *Vec) Clear() {
	v.oldHead = nil
	v.newTail = v.newTail[:0]
}

// Contains returns true if the element is inside the array
func (v *Vec) Contains(el Element) bool {
	bigger := v.newTail
	smaller := v.oldHead
	if len(v.newTail) < v.oldLen() {
		bigger, smaller = smaller, bigger
	}
	for i := range bigger {
		if bigger[i] == el {
			return true
		}
		if i < len(smaller) && smaller[i] == el {
			return true
		}
	}
	return false
}

// ContainsCmp returns true if the element is inside the array, will use the cmp func
func (v *Vec) ContainsCmp(el Element, cmp func(arrayElement Element, el Element) bool) bool {
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
func (v *Vec) First() Element {
	if v.oldLen() > 0 {
		return v.oldHead[0]
	}
	if len(v.newTail) > 0 {
		return v.newTail[0]
	}
	return nil
}

// Last returns the last element of the array, returns null if it is empty
func (v *Vec) Last() Element {
	if len(v.newTail) > 0 {
		return v.newTail[len(v.newTail)-1]
	}
	oldLen := v.oldLen()
	if oldLen > 0 {
		return v.oldHead[oldLen]
	}
	return nil
}

// PopFront pops the first element of the array, returns null if the array is empty
func (v *Vec) PopFront() Element {
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
	return nil
}

// PopBack pops the last element of the array, returns null if the array is empty
func (v *Vec) PopBack() Element {
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
	return nil
}

// Pop same as PopBack
func (v *Vec) Pop() Element {
	return v.PopBack()
}

// Iter generates an array of elements (allocates space for the iteration)
func (v *Vec) Iter() []Element {
	elements := make([]Element, 0, v.Len())
	if v.oldHead != nil {
		elements = append(elements, v.oldHead...)
	}
	elements = append(elements, v.newTail...)
	return elements
}

// IterFunc iterates through the array doing a callback to the passed function
func (v *Vec) IterFunc(fn func(el Element, index int)) {
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
func (v *Vec) Push(el Element) {
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
func (v *Vec) Append(el ...Element) {
	for _, e := range el {
		v.Push(e)
	}
}

func (v *Vec) carry() {
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

const pushMultiplierOldVector = 2

func (v *Vec) grow(growFactor int) {
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
	//
	// We also need to make sure we can fit the additional capacity required for `extra`.
	// Normally, that'll be handled by `pushes`, but not always!
	add := max(pushes, growFactor)
	elements := make([]Element, 0, add+pushes+need)

	v.oldHead = make([]Element, 0, add+pushes+need*pushMultiplierOldVector)
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
func (v *Vec) Debug() string {
	return fmt.Sprintf("Old: %v \n", v.oldHead) + fmt.Sprintf("New: %v", v.newTail)
}
