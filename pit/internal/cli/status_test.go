package cli

import "testing"

func TestStatusNeverSigns(t *testing.T) {
	if !StatusNeverSigns("session   none on this CLI until desktop or keychain bind") {
		t.Fatal("status")
	}
	if StatusNeverSigns("session_key=abc") {
		t.Fatal("secret")
	}
}
