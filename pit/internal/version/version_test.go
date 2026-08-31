package version

import "testing"

func TestString(t *testing.T) {
	if String() != "PIT 0.9.6" {
		t.Fatal(String())
	}
}
