package sdk

import "testing"

func TestSDKCannotPostExchange(t *testing.T) {
	if (Client{}).CanPostExchange() {
		t.Fatal("sdk")
	}
}
