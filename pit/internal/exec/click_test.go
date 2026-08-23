package exec

import "testing"

func TestDuplicateClick(t *testing.T) {
	if err := DuplicateClick(true, nil); err != nil {
		t.Fatal(err)
	}
	if err := DuplicateClick(false, nil); err == nil {
		t.Fatal("dup")
	}
}
