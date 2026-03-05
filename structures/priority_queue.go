package structures

import "container/heap"

// PriorityItem represents an element in the priority queue
type PriorityItem[T any] struct {
	value    T   // The actual data you want to store
	priority int // Priority of the item
	index    int // Index in the heap (maintained by heap.Interface)
}

func (item *PriorityItem[T]) Value() T {
	return item.value
}

// NewPriorityItem creates a new PriorityItem
func NewPriorityItem[T any](value T, priority int) *PriorityItem[T] {
	return &PriorityItem[T]{value: value, priority: priority}
}

// PriorityQueue implements heap.Interface
type PriorityQueue[T any] []*PriorityItem[T]

// Len returns the number of elements in the queue
func (pq *PriorityQueue[T]) Len() int { return len(*pq) }

// Less defines the priority order
// For min-heap (the smallest priority first): pq[i].priority < pq[j].priority
// For max-heap (the largest priority first): pq[i].priority > pq[j].priority
func (pq *PriorityQueue[T]) Less(i, j int) bool {
	return (*pq)[i].priority < (*pq)[j].priority // Min-heap
}

// Swap swaps the elements with indexes i and j.
func (pq *PriorityQueue[T]) Swap(i, j int) {
	(*pq)[i], (*pq)[j] = (*pq)[j], (*pq)[i]
	(*pq)[i].index = i
	(*pq)[j].index = j
}

func (pq *PriorityQueue[T]) Push(x any) {
	item := x.(*PriorityItem[T])
	item.index = len(*pq)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue[T]) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// Update modifies the priority of an item in the queue
func (pq *PriorityQueue[T]) Update(item *PriorityItem[T], value T, priority int) {
	item.value = value
	item.priority = priority
	heap.Fix(pq, item.index)
}
