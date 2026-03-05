package cmd

import (
	"cmp"
	"container/heap"

	"github.com/apex/log"
	"github.com/spf13/cobra"
	"gitlab.com/mshindle/tidbits/structures"
)

// pqCmd represents the pq command
var pqCmd = &cobra.Command{
	Use:   "pq",
	Short: "use PriorityQueue to solve a series of common problems",
	Long: `
Solve some common problems using a priority queue. This command demonstrates how to use a priority queue to efficiently
solve problems that involve prioritizing tasks or events based on their importance or urgency.

+-------------------+-------------------------------+----------------------------------+
|   Pattern         |   Example Problems            |   Why PQ Works                   |
+-------------------+-------------------------------+----------------------------------+
| Kth Something     | Kth largest, Top K frequent   | Keep only K elements, O(n log k) |
| Merge Sorted      | Merge K lists, Smallest range | Always pick smallest/largest     |
| Running Stats     | Median stream, Sliding window | Maintain dynamic order fast      |
| Graph Algorithms  | Dijkstra, Prim, A*            | Expand lowest-cost node first    |
| Scheduling        | CPU tasks, Meeting rooms      | Process by priority/availability |
+-------------------+-------------------------------+----------------------------------+
`,
	RunE: pq,
}

func init() {
	rootCmd.AddCommand(pqCmd)
	//pqCmd.Flags().Int64VarP(&points, "points", "p", 50000000, "number of points to use for calculation")
	//pqCmd.Flags().IntVarP(&numWorkers, "workers", "w", 3, "number of workers to calculate if points are in or out of circle")
}

func pq(cmd *cobra.Command, args []string) error {
	// kth largest
	nums := []int{71, 12, 33, 54, 65, 16, 47, 28, 89, 10}
	k := 4
	log.WithField("nums", nums).WithField("k", k).Info("pq kth largest started")
	kth := pqKthLargestElement(nums, k)
	log.WithField("largest", kth).Info("pq kth largest calculated")

	// merge k lists
	nodeLists := make([]*structures.ListNode[int], 3)
	for i := 0; i < 8; i = i + 3 {
		n := &structures.ListNode[int]{
			Value: i,
			Next: &structures.ListNode[int]{
				Value: i + 1,
				Next: &structures.ListNode[int]{
					Value: i + 2,
				},
			},
		}
		nodeLists = append(nodeLists, n)
	}
	mergedNode := pqMergeKLists[int](nodeLists)
	for mergedNode != nil {
		log.WithField("value", mergedNode.Value).Info("merged node")
		mergedNode = mergedNode.Next
	}

	return nil
}

// pqKthLargestElement finds the kth largest element in a given array.
// How this works: We maintain a min-heap of the K largest elements we’ve seen. When we have more than K elements,
// we kick out the smallest one. In the end, the smallest element in our heap of K largest elements is the
// Kth largest overall.
//
// Time complexity: O(n log k) instead of O(n log n) for sorting. If k is small, this is a huge win.
func pqKthLargestElement(nums []int, k int) int {
	pq := make(structures.PriorityQueue[int], 0)
	heap.Init(&pq)

	for _, num := range nums {
		item := structures.NewPriorityItem(num, num)
		heap.Push(&pq, item)
		if pq.Len() > k {
			heap.Pop(&pq)
		}
	}
	return pq[0].Value()
}

type listNodeHeap[T cmp.Ordered] []*structures.ListNode[T]

func (h listNodeHeap[T]) Len() int           { return len(h) }
func (h listNodeHeap[T]) Less(i, j int) bool { return h[i].Value < h[j].Value }
func (h listNodeHeap[T]) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *listNodeHeap[T]) Push(x any) {
	*h = append(*h, x.(*structures.ListNode[T]))
}
func (h *listNodeHeap[T]) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func pqMergeKLists[T cmp.Ordered](lists []*structures.ListNode[T]) *structures.ListNode[T] {
	h := &listNodeHeap[T]{}
	heap.Init(h)

	for _, node := range lists {
		if node != nil {
			log.WithField("value", node.Value).Info("lists read node")
			heap.Push(h, node)
		}
	}

	dummy := &structures.ListNode[T]{Next: nil}
	current := dummy
	for h.Len() > 0 {
		node := heap.Pop(h).(*structures.ListNode[T])
		current.Next = node
		current = current.Next
		if node.Next != nil {
			heap.Push(h, node.Next)
		}
	}

	return dummy.Next
}
