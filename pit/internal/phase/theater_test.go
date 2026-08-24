package phase

import "testing"

func TestRefuseTheater(t *testing.T) {
	if err := RefuseTheater(WaitingForUser); err != nil {
		t.Fatal(err)
	}
	if err := RefuseTheater("SPINNING"); err == nil {
		t.Fatal("theater")
	}
}
