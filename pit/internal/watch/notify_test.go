package watch

import "testing"

func TestMayPlaceOrder(t *testing.T) {
	if err := MayPlaceOrder(true); err == nil {
		t.Fatal("watch")
	}
	if err := MayPlaceOrder(false); err != nil {
		t.Fatal(err)
	}
}
