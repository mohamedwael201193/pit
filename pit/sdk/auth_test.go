package sdk

import "testing"

func TestClientCannotReadAuthFile(t *testing.T) {
	if (Client{}).CanReadAuthFile() {
		t.Fatal("auth")
	}
}
