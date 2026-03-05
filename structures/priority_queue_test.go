package structures

import (
	"container/heap"
	"testing"
)

func TestPriorityQueue_Basic(t *testing.T) {
	pq := &PriorityQueue[string]{}
	heap.Init(pq)

	items := []*PriorityItem[string]{
		{value: "orange", priority: 3},
		{value: "apple", priority: 1},
		{value: "banana", priority: 2},
	}

	for _, item := range items {
		heap.Push(pq, item)
	}

	if pq.Len() != 3 {
		t.Errorf("expected Len 3, got %d", pq.Len())
	}

	expected := []string{"apple", "banana", "orange"}
	for _, val := range expected {
		item := heap.Pop(pq).(*PriorityItem[string])
		if item.value != val {
			t.Errorf("expected %s, got %s", val, item.value)
		}
	}

	if pq.Len() != 0 {
		t.Errorf("expected Len 0, got %d", pq.Len())
	}
}

func TestPriorityQueue_UpdateReceiver(t *testing.T) {
	pq := &PriorityQueue[string]{}
	heap.Init(pq)

	item := &PriorityItem[string]{value: "initial", priority: 10}
	heap.Push(pq, item)
	heap.Push(pq, &PriorityItem[string]{value: "other", priority: 5})

	// Use the Update receiver
	pq.Update(item, "updated", 1)

	if item.value != "updated" {
		t.Errorf("expected updated, got %s", item.value)
	}

	popped := heap.Pop(pq).(*PriorityItem[string])
	if popped.value != "updated" {
		t.Errorf("expected updated, got %s", popped.value)
	}
	if popped.priority != 1 {
		t.Errorf("expected priority 1, got %d", popped.priority)
	}
}

func TestPriorityQueue_Update(t *testing.T) {
	pq := &PriorityQueue[string]{}
	heap.Init(pq)

	item := &PriorityItem[string]{value: "initial", priority: 10}
	heap.Push(pq, item)
	heap.Push(pq, &PriorityItem[string]{value: "other", priority: 5})

	// Update the priority of "initial" to be highest
	item.priority = 1
	heap.Fix(pq, item.index)

	popped := heap.Pop(pq).(*PriorityItem[string])
	if popped.value != "initial" {
		t.Errorf("expected initial, got %s", popped.value)
	}
}

func TestPriorityQueue_Remove(t *testing.T) {
	pq := &PriorityQueue[string]{}
	heap.Init(pq)

	item1 := &PriorityItem[string]{value: "item1", priority: 1}
	item2 := &PriorityItem[string]{value: "item2", priority: 2}
	item3 := &PriorityItem[string]{value: "item3", priority: 3}

	heap.Push(pq, item1)
	heap.Push(pq, item2)
	heap.Push(pq, item3)

	// Remove item2 from the middle
	removed := heap.Remove(pq, item2.index).(*PriorityItem[string])
	if removed.value != "item2" {
		t.Errorf("expected item2, got %s", removed.value)
	}

	if pq.Len() != 2 {
		t.Errorf("expected Len 2, got %d", pq.Len())
	}

	// Check remaining items
	if p1 := heap.Pop(pq).(*PriorityItem[string]); p1.value != "item1" {
		t.Errorf("expected item1, got %s", p1.value)
	}
	if p3 := heap.Pop(pq).(*PriorityItem[string]); p3.value != "item3" {
		t.Errorf("expected item3, got %s", p3.value)
	}
}

func TestPriorityQueue_EmptyPop(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic on Pop from empty queue")
		}
	}()

	pq := &PriorityQueue[string]{}
	heap.Init(pq)
	heap.Pop(pq)
}

func TestPriorityQueue_SamePriority(t *testing.T) {
	pq := &PriorityQueue[string]{}
	heap.Init(pq)

	heap.Push(pq, &PriorityItem[string]{value: "a", priority: 1})
	heap.Push(pq, &PriorityItem[string]{value: "b", priority: 1})

	if pq.Len() != 2 {
		t.Errorf("expected Len 2, got %d", pq.Len())
	}

	// The order for the same priority is not strictly guaranteed by heap,
	// but it should at least not crash and return both.
	heap.Pop(pq)
	heap.Pop(pq)

	if pq.Len() != 0 {
		t.Errorf("expected Len 0, got %d", pq.Len())
	}
}
