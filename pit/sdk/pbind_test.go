package sdk

import "testing"

func TestSDKCannotAuthorizePreview(t *testing.T) {
	if (Client{}).CanAuthorizePreview() {
		t.Fatal("sdk")
	}
}
