/*
* @Author: souravray
* @Date:   2014-11-01 23:40:31
* @Last Modified by:   souravray
* @Last Modified time: 2014-11-02 00:09:05
 */

package heap

import (
	"container/heap"
)

// extend container/heap.Interface to implement PopWithoutWait method in a PriorityWaitQueue metod
type Interface interface {
	heap.Interface
	PopWithoutWait() interface{}
}

// Init calls container/heap.Init methode
func Init(h Interface) {
	heap.Init(h)
}

// Push calls container/heap.Push methode
func Push(h Interface, x interface{}) {
	heap.Push(h, x)
}

// Pop calls container/heap.Pop methode
func Pop(h Interface) interface{} {
	return heap.Pop(h)
}

// Fix calls container/heap.Fix methode
func Fix(h Interface, i int) {
	heap.Fix(h, i)
}

// Remove removes the element at index i from the heap.
func Remove(h Interface, i int) interface{} {
	n := h.Len() - 1
	if n != i {
		h.Swap(i, n)
		down(h, i, n)
		up(h, i)
	}
	return h.PopWithoutWait()
}

// private up and down mwthods
func up(h Interface, j int) {
	for {
		i := (j - 1) / 2 // parent
		if i == j || !h.Less(j, i) {
			break
		}
		h.Swap(i, j)
		j = i
	}
}

func down(h Interface, i, n int) {
	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && !h.Less(j1, j2) {
			j = j2 // = 2*i + 2  // right child
		}
		if !h.Less(j, i) {
			break
		}
		h.Swap(i, j)
		i = j
	}
}
