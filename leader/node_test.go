package leader

import "testing"

func TestNode_IsLowerThan(t *testing.T) {
	// Given
	node := Node{ID: "node-02", Rank: 2}

	// Then
	if !node.IsRankHigher(1) {
		t.Errorf("Rank 2 should be higher than 1")
	}
}
