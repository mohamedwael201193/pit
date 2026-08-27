package exec

import "testing"

func TestNeedOnVenue(t *testing.T) {
	if err := NeedOnVenue(false); err == nil {
		t.Fatal("missing")
	}
	if err := NeedOnVenue(true); err != nil {
		t.Fatal(err)
	}
}
