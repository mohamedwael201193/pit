package sdk

import "testing"

func TestSurfaceNeverSignsAndEmptyWatch(t *testing.T) {
	c := Client{}
	if c.CanHoldSession() {
		t.Fatal("session")
	}
	if c.EmptyWatch() == "" {
		t.Fatal("watch")
	}
	if len(c.Phases()) < 8 {
		t.Fatal("phases")
	}
}
