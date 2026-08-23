package memory

import (
	"testing"

	"github.com/mohamedwael201193/pit/internal/identity"
)

func TestListDenied(t *testing.T) {
	a := identity.NewWorkspaceID()
	b := identity.NewWorkspaceID()
	if err := ListDenied(b, a, "thesis", "note"); err != nil {
		t.Fatal(err)
	}
	if err := ListDenied(a, a, "thesis", "note"); err == nil {
		t.Fatal("same")
	}
}
