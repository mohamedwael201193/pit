package watch

import "testing"

func TestWatchCannotPlace(t *testing.T) {
	if err := MayPlaceOrder(true); err == nil {
		t.Fatal("watch")
	}
	if err := MayPlaceOrder(false); err != nil {
		t.Fatal(err)
	}
}

func TestNotifyEmptyIsReal(t *testing.T) {
	if NotifyCopy(0) != "No opportunities match your policy." {
		t.Fatal(NotifyCopy(0))
	}
}
