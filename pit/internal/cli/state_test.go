package cli

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/mohamedwael201193/pit/internal/identity"
)

func TestSaveLoadKill(t *testing.T) {
	dir := t.TempDir()
	id := identity.NewWorkspaceID()
	st := DiskState{WorkspaceID: id, Network: "testnet", Wallet: "0x1111111111111111111111111111111111111111"}
	if err := Save(dir, st); err != nil {
		t.Fatal(err)
	}
	got, err := Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	if got.WorkspaceID != id || got.Kill {
		t.Fatalf("%+v", got)
	}
	if err := SetKill(dir, true); err != nil {
		t.Fatal(err)
	}
	got, _ = Load(dir)
	if !got.Kill {
		t.Fatal("kill")
	}
	if runtime.GOOS == "windows" {
		return
	}
	mode, err := os.Stat(filepath.Join(dir, "state.json"))
	if err != nil {
		t.Fatal(err)
	}
	if mode.Mode().Perm()&0o077 != 0 {
		t.Fatalf("perms %o", mode.Mode().Perm())
	}
}

func TestSaveRejectsBadWorkspace(t *testing.T) {
	if err := Save(t.TempDir(), DiskState{WorkspaceID: "nope", Network: "mainnet", Wallet: "0x1111111111111111111111111111111111111111"}); err == nil {
		t.Fatal("id")
	}
}
