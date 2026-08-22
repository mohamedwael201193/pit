package sdk

import (
	"testing"

	"github.com/mohamedwael201193/pit/internal/engine"
	"github.com/mohamedwael201193/pit/internal/identity"
)

func TestForecastForHidesOtherWorkspace(t *testing.T) {
	a := identity.NewWorkspaceID()
	b := identity.NewWorkspaceID()
	f, err := engine.BuildForecast("hyperliquid:perp:ETH", "buy", "mark below 2000", 0.2, 0.55)
	if err != nil {
		t.Fatal(err)
	}
	_, err = Client{}.ForecastFor(b, engine.StoredForecast{Workspace: a, ID: "f1", Forecast: f})
	if err == nil {
		t.Fatal("idor")
	}
}
