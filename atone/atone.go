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

// Element interface
type Element interface{}

// Vec is the implementation of an atone vector just like in https://github.com/jonhoo/atone, there you will find what is so special about this implementation
type Vec struct {
	oldHead []Element
	newTail []Element
}

// NItemsToMoveOnEachInsert is the number of items we move on each insert, recommended values: 1, 4, 6....
const NItemsToMoveOnEachInsert = 5

// New returns a new atone Vec
func New() *Vec {
	return &Vec{
		newTail: make([]Element, 0, 0),
		oldHead: nil,
	}
}

// From returns a new Vec from a Slice
func From(elements []Element) *Vec {
	return &Vec{
		newTail: elements,
		oldHead: nil,
	}
}

// NewWithCapacity is the equivalent of doing make([]Element, 0, capacity)
func NewWithCapacity(capacity uint64) *Vec {
	return &Vec{
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

// Lookup returns an element, the boolean is false if the element does not exist.
func (v *Vec) Lookup(index int) (Element, bool) {
	if index < v.oldLen() {
		return v.oldHead[index], true
	}
	offset := index - v.oldLen()
	if offset >= len(v.newTail) || offset < 0 {
		return nil, false
	}
	return v.newTail[offset], true
}

// Get returns the element in the specified index, can panic if it is outofbounds, if you don't want to panic on get, use Lookup
func (v *Vec) Get(index int) Element {
	if index < v.oldLen() {
		return v.oldHead[index]
	}

	offset := index - v.oldLen()
	return v.newTail[offset]
}

// GetRef returns FOR SURE a pointer to the element even though it is a stack element like int
func (v *Vec) GetRef(index int) *Element {
	if index < v.oldLen() {
		return &v.oldHead[index]
	}

	offset := index - v.oldLen()
	return &v.newTail[offset]
}

// Find02 tries to find not doing a continuous loop
func (v *Vec) Find02(el Element) int {
	bigger := v.newTail
	smaller := v.oldHead
	if len(v.newTail) < v.oldLen() {
		bigger, smaller = smaller, bigger
	}
	for i := range bigger {
		if bigger[i] == el {
			return i
		}
		if i < len(smaller) && smaller[i] == el {
			return i
		}
	}
	return -1
}

// Find finds doing a lookup in head and then in tail
func (v *Vec) Find(el Element) int {
	if v.oldHead != nil {
		for i := range v.oldHead {
			if v.oldHead[i] == el {
				return i
			}
		}
	}

	for i := range v.newTail {
		if v.newTail[i] == el {
			return i
		}
	}

	return -1
}

func find(els []Element, el Element) int {
	for i := range els {
		if els[i] == el {
			return i
		}
	}
	return -1
}

// FindMultithreaded finds an element with multithreading (with a lot of elements 1000000+)
func (v *Vec) FindMultithreaded(el Element) int {
	channel := make(chan int, 2)
	go func() {
		channel <- find(v.oldHead, el)
	}()
	go func() {
		channel <- find(v.newTail, el)
	}()
	for found := range channel {
		if found != -1 {
			return found
		}
	}
	return -1
}

// Insert .
func (v *Vec) Insert(el Element) {

	if len(v.newTail) == cap(v.newTail) {
		v.grow(1)
		v.Insert(el)
		return
	}
	if v.oldLen() == 0 {
		els := make([]Element, 0, cap(v.newTail)+1)
		els = append(els, el)
		v.newTail = append(els, v.newTail...)
		return
	}
	// storage for sufficient elements in the new tail. maybe with a better implementaiton we could jump this
	els := make([]Element, 0, cap(v.newTail)+1)
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
func (v *Vec) Swap(i int, j int) {
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

func reverseSlice(s []Element) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// Reverse inplace the array and empties the old head
func (v *Vec) Reverse() {
	reverseSlice(v.newTail)
	if v.oldHead != nil {
		for i := range v.oldHead {
			v.newTail = append(v.newTail, v.oldHead[v.oldLen()-i-1])
		}
		v.oldHead = nil
	}
}

// Reserve the desired size inmemory to let space for nElements, it might reserve more memory than necessary for leaving space for more items for carry()
func (v *Vec) Reserve(nElements int) {

	if v.oldLen() > 0 {

		v.carryAll()
		v.grow(nElements)
		return
	}
	v.grow(nElements)
}

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
func (v *Vec) Truncate(n int) {
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

// Slice generates a slice slicing the array from start to end (end is not inclusive and start is)
func (v *Vec) Slice(start, end int) []Element {
	elements := make([]Element, 0, end-start)
	if v.oldLen() > start {
		newEnd := end
		if newEnd > v.oldLen() {
			newEnd = v.oldLen()
		}
		elements = append(elements, v.oldHead[start:newEnd]...)
		if newEnd > v.oldLen() {
			elements = append(elements, v.newTail[:newEnd-v.oldLen()])
		}
		return elements
	}
	elements = append(elements, v.newTail[start:end]...)
	return elements
}

// Array creates a slice of this array
func (v *Vec) Array() []Element {
	return v.Slice(0, v.Len())
}

// SliceThis returns a new Vec with the specified slice (end non inclusive and start is inclusive)
func (v *Vec) SliceThis(start, end int) *Vec {
	return From(v.Slice(start, end))
}

// ForEach iterates through the array doing a callback to the passed function
func (v *Vec) ForEach(fn func(el Element, index int)) {
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

func (v *Vec) carryAll() {
	if v.oldLen() == 0 {
		v.oldHead = nil
		return
	}
	v.newTail = append(v.oldHead[0:], v.newTail...)
	v.oldHead = nil
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
	// We also need to make sure we can fit the additional capacity required for `extra`.
	// Normally, that'll be handled by `pushes`, but not always!
	add := max(pushes, growFactor)
	elements := make([]Element, 0, cap(v.newTail)+add+pushes+need)
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
