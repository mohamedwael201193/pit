package version

import "testing"

func TestString(t *testing.T) {
	if String() != "PIT 0.2.2" {
		t.Fatal(String())
	}
}
