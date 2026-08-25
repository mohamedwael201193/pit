package calib

import "testing"

func TestNeedResolved(t *testing.T) {
	if NeedResolved() != 30 {
		t.Fatal(NeedResolved())
	}
	if err := RefuseSparse(0); err != nil {
		t.Fatal(err)
	}
	if err := RefuseSparse(29); err != nil {
		t.Fatal(err)
	}
}
