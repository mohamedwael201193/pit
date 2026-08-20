package watch

import "testing"

func TestNoCloudWatchAndNoGhost(t *testing.T) {
	if err := CloudLoopForbidden(); err == nil {
		t.Fatal("cloud")
	}
	if err := GhostCard(-1); err == nil {
		t.Fatal("ghost")
	}
	if err := GhostCard(0); err != nil {
		t.Fatal(err)
	}
}
