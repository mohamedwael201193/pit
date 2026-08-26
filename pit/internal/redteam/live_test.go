package redteam

import "testing"

func TestExpiredAuthorize(t *testing.T) {
	if err := ExpiredAuthorize(); err == nil {
		t.Fatal("empty session")
	}
}
