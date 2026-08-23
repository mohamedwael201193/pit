package watch

import (
	"testing"

	"github.com/mohamedwael201193/pit/internal/identity"
)

func TestBound(t *testing.T) {
	if err := Bound(identity.NewWorkspaceID()); err != nil {
		t.Fatal(err)
	}
	if err := Bound("not-a-workspace"); err == nil {
		t.Fatal("id")
	}
}
