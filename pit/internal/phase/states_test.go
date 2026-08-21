package phase

import "testing"

func TestKnownStates(t *testing.T) {
	if !Known(WaitingForUser) || Known("LOADING_SPINNER") {
		t.Fatal("states")
	}
}
