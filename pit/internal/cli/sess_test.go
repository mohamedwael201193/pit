package cli

import (
	"strings"
	"testing"

	"github.com/mohamedwael201193/pit/internal/identity"
)

func TestSaveLoadSessionMeta(t *testing.T) {
	dir := t.TempDir()
	ws := identity.NewWorkspaceID()
	in := SessionFile{
		ID:        identity.NewWorkspaceID(),
		AgentAddr: "0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
		Workspace: ws,
		Network:   "testnet",
		PolicyVer: "v1",
		Expires:   50,
	}
	if err := SaveSession(dir, in); err != nil {
		t.Fatal(err)
	}
	got, err := LoadSession(dir)
	if err != nil {
		t.Fatal(err)
	}
	if got.AgentAddr != in.AgentAddr || got.Workspace != ws {
		t.Fatalf("%+v", got)
	}
}

func TestRefuseSessionSecret(t *testing.T) {
	if err := RefuseSessionSecret([]byte(`{"session_key":"nope"}`)); err == nil {
		t.Fatal("secret")
	}
	if err := SaveSession(t.TempDir(), SessionFile{}); err == nil {
		t.Fatal("empty")
	}
}

func TestLoadMissingSession(t *testing.T) {
	if _, err := LoadSession(t.TempDir()); err == nil || !strings.Contains(err.Error(), "session_expired") {
		t.Fatalf("%v", err)
	}
}
