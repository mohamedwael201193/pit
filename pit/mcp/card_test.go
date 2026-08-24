package mcp

import "testing"

func TestEmptyCard(t *testing.T) {
	r := EmptyCard()
	body, _ := r.Body.(map[string]any)
	if body["invented"] != false || body["n"] != 0 {
		t.Fatal(body)
	}
}
