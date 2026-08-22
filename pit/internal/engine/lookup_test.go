package engine

import (
	"testing"

	"github.com/mohamedwael201193/pit/internal/identity"
)

func TestLookupHidesOtherWorkspace(t *testing.T) {
	a := identity.NewWorkspaceID()
	b := identity.NewWorkspaceID()
	f, err := BuildForecast("hyperliquid:perp:ETH", "buy", "mark below 2000", 0.2, 0.55)
	if err != nil {
		t.Fatal(err)
	}
	rec := StoredForecast{Workspace: a, ID: "f1", Forecast: f}
	if _, err := Lookup(b, rec); err == nil {
		t.Fatal("idor")
	}
	got, err := Lookup(a, rec)
	if err != nil || got.Invalidation == "" {
		t.Fatal(err)
	}
}
