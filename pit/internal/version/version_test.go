package version

import "testing"

func TestString(t *testing.T) {
	if String() != "PIT 0.1.5" {
		t.Fatal(String())
	}
}
