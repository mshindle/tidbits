package concurrency

import (
	"testing"
)

func TestLeaking(t *testing.T) {
	if !Leaking() {
		t.Error("expected Leaking() to return true, but got false")
	}
}

func TestLeakingCascade(t *testing.T) {
	if !LeakingCascade() {
		t.Error("expected LeakingCascade() to return true, but got false")
	}
}
