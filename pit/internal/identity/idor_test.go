package identity

import "testing"

func TestSameWorkspace(t *testing.T) {
	a := NewWorkspaceID()
	b := NewWorkspaceID()
	if err := SameWorkspace(a, a); err != nil {
		t.Fatal(err)
	}
	if err := SameWorkspace(a, b); err == nil {
		t.Fatal("idor")
	}
}
