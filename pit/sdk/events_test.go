package sdk

import "testing"

func TestEventsAreNamedStates(t *testing.T) {
	c := Client{}
	if len(c.Events()) < 10 {
		t.Fatal("events")
	}
	for _, e := range c.Events() {
		if !c.EventKnown(e) {
			t.Fatal(e)
		}
	}
	if c.EventKnown("SPINNING") {
		t.Fatal("theater")
	}
}
