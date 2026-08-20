package identity

import "testing"

func TestNormalizeAddress(t *testing.T) {
	a, err := NormalizeAddress("0xBDfCeE82Bd42FEfA58ee850B3709636a8B6b0034")
	if err != nil {
		t.Fatal(err)
	}
	if a != "0xbdfcee82bd42fefa58ee850b3709636a8b6b0034" {
		t.Fatalf("got %s", a)
	}
	if _, err := NormalizeAddress("BDfC"); err == nil {
		t.Fatal("expected error")
	}
	if _, err := NormalizeAddress("0xzz"); err == nil {
		t.Fatal("expected error")
	}
}

func TestParseWorkspaceIDHidesFormat(t *testing.T) {
	id := NewWorkspaceID()
	got, err := ParseWorkspaceID(id)
	if err != nil || got != id {
		t.Fatalf("%s %v", got, err)
	}
	if _, err := ParseWorkspaceID("not-a-uuid"); err == nil || err.Error() != "not found" {
		t.Fatalf("guess must return not found, got %v", err)
	}
}
