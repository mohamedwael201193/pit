package cli

import "testing"

func TestExpiredCopy(t *testing.T) {
	if ExpiredCopy == "" || RevokedCopy == "" {
		t.Fatal("copy")
	}
}
