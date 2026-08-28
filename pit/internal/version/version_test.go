package version

import "testing"

func TestString(t *testing.T) {
	if String() != "PIT 0.2.0" {
		t.Fatal(String())
	}
}
