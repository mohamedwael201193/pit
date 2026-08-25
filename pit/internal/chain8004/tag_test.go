package chain8004

import "testing"

func TestRefusePrivateTag(t *testing.T) {
	if err := RefusePrivateTag("position-leak"); err == nil {
		t.Fatal("pos")
	}
	if err := RefusePrivateTag("resolved"); err != nil {
		t.Fatal(err)
	}
}
