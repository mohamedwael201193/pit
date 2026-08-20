package cli

import "testing"

func TestPipedYesDenied(t *testing.T) {
	if err := ConfirmAuthorize(false, "yes", true); err == nil {
		t.Fatal("piped")
	}
	if err := ConfirmAuthorize(true, "yes", true); err == nil {
		t.Fatal("yes is not AUTHORIZE")
	}
	if err := ConfirmAuthorize(true, "AUTHORIZE", false); err == nil {
		t.Fatal("flag")
	}
	if err := ConfirmAuthorize(true, "AUTHORIZE", true); err != nil {
		t.Fatal(err)
	}
}
