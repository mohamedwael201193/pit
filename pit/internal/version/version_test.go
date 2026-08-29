package version

import "testing"

func TestString(t *testing.T) {
	if String() != "PIT 0.4.0" {
		t.Fatal(String())
	}
}
