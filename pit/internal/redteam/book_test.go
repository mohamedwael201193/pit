package redteam

import "testing"

func TestBookNeverInHealth(t *testing.T) {
	if !BookNeverInHealth(map[string]any{"ok": true, "sign": false, "trade": false}) {
		t.Fatal("clean")
	}
	if BookNeverInHealth(map[string]any{"private_book": []byte("x")}) {
		t.Fatal("book")
	}
}
