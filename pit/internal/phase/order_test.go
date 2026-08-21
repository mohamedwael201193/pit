package phase

import "testing"

func TestNextStopsAtFailed(t *testing.T) {
	if _, err := Next(Failed); err == nil {
		t.Fatal("failed")
	}
	n, err := Next(Connecting)
	if err != nil || n != Authenticating {
		t.Fatalf("%s %v", n, err)
	}
	if !Known(n) {
		t.Fatal(n)
	}
}
