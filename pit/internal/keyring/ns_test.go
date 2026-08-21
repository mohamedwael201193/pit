package keyring

import (
	"strings"
	"testing"

	"github.com/mohamedwael201193/pit/internal/identity"
)

func TestNamespaceIsolatesWorkspaces(t *testing.T) {
	a := identity.NewWorkspaceID()
	b := identity.NewWorkspaceID()
	na, err := Namespace("mainnet", a, "session")
	if err != nil {
		t.Fatal(err)
	}
	nb, err := Namespace("mainnet", b, "session")
	if err != nil {
		t.Fatal(err)
	}
	if na == nb || !strings.Contains(na, a) {
		t.Fatalf("%s %s", na, nb)
	}
	if _, err := Namespace("mainnet", "not-a-uuid", "session"); err == nil {
		t.Fatal("id")
	}
}
