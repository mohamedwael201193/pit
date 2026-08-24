package cli

import "testing"

func TestRefusePrint(t *testing.T) {
	if err := RefusePrint("workspace ready"); err != nil {
		t.Fatal(err)
	}
	if err := RefusePrint("session_key=0xabc"); err == nil {
		t.Fatal("print")
	}
}
