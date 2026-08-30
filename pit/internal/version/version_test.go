package version

import "testing"

func TestString(t *testing.T) {
	if String() != "PIT 0.9.4" {
		t.Fatal(String())
	}
}
