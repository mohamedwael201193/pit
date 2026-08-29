package version

import "testing"

func TestString(t *testing.T) {
	if String() != "PIT 0.6.1" {
		t.Fatal(String())
	}
}
