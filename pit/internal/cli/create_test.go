package cli

import (
	"testing"

	"github.com/mohamedwael201193/pit/internal/identity"
	"github.com/mohamedwael201193/pit/internal/session"
)

func TestCreateLocalSessionDoesNotExport(t *testing.T) {
	dir := t.TempDir()
	ws := identity.NewWorkspaceID()
	sf, err := CreateLocalSession(dir, ws, "mainnet", "v1")
	if err != nil {
		t.Fatal(err)
	}
	if sf.AgentAddr == "" || sf.Expires == 0 {
		t.Fatalf("%+v", sf)
	}
	s, err := LiveFromDisk(dir, false, sf.Expires-1)
	if err != nil {
		t.Fatal(err)
	}
	if s.Workspace != ws {
		t.Fatal(s.Workspace)
	}
	if _, err := LiveFromDisk(dir, false, sf.Expires); err == nil {
		t.Fatal("expired")
	}
	if session.Alive(s, sf.Expires) {
		t.Fatal("boundary")
	}
}
