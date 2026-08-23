package compute

import "testing"

func TestTamperFails(t *testing.T) {
	ok := SchemeE2EE + ":body"
	if err := TamperFails(ok, ok); err != nil {
		t.Fatal(err)
	}
	if err := TamperFails(ok, SchemeE2EE+":mutated"); err == nil {
		t.Fatal("tamper")
	}
	if err := TamperFails(ok, "plaintext"); err == nil {
		t.Fatal("scheme")
	}
}
