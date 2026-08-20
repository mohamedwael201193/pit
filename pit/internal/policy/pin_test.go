package policy

import "testing"

func TestPinRoundTrip(t *testing.T) {
	p := Default()
	h, err := p.Hash()
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	if _, err := PinFile(dir, "ws-a", h); err != nil {
		t.Fatal(err)
	}
	got, err := ReadPin(dir, "ws-a")
	if err != nil || got != h {
		t.Fatal(got, err)
	}
	if _, err := ReadPin(dir, "ws-b"); err == nil {
		t.Fatal("cross workspace")
	}
}
