package cli

import "testing"

func TestStatusNeverSigns(t *testing.T) {
	if !StatusNeverSigns("session   none on this CLI until desktop or keychain bind") {
		t.Fatal("status")
	}
	if StatusNeverSigns("session_key=abc") {
		t.Fatal("secret")
	}
	if !StatusNeverSigns(LinkCopy(false, nil)) || !StatusNeverSigns(VenueCopy(true, nil)) {
		t.Fatal("copy")
	}
}
