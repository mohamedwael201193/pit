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

func TestGuessedWorkspaceDenied(t *testing.T) {
	if err := SameWorkspace("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", NewWorkspaceID()); err == nil {
		t.Fatal("guessed id")
	}
	if err := SameWorkspace("not-an-id", NewWorkspaceID()); err == nil {
		t.Fatal("garbage id")
	}
}
