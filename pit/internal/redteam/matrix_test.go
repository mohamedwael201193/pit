package redteam

import "testing"

func TestS3Allowlist(t *testing.T) {
	if errs := SessionActionMatrix(); len(errs) != 0 {
		t.Fatalf("%v", errs)
	}
}
