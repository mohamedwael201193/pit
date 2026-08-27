package compute

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDiscoverSealer(t *testing.T) {
	dir := t.TempDir()
	if DiscoverSealer(dir) != "" {
		t.Fatal("empty")
	}
	p := filepath.Join(dir, "pit-sealer")
	if err := os.WriteFile(p, []byte("x"), 0o700); err != nil {
		t.Fatal(err)
	}
	got := DiscoverSealer(dir)
	if got != p {
		t.Fatalf("%s", got)
	}
}

func TestDiscoverSealerWindowsSidecarName(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "pit-sealer-x86_64-pc-windows-msvc.exe")
	if err := os.WriteFile(p, []byte("x"), 0o700); err != nil {
		t.Fatal(err)
	}
	got := DiscoverSealer(dir)
	if got != p {
		t.Fatalf("%s", got)
	}
}

func TestLookBinPrefersEnv(t *testing.T) {
	t.Setenv("PIT_COMMITTEE_BIN", "/opt/pit-sealer")
	if LookBin() != "/opt/pit-sealer" {
		t.Fatal(LookBin())
	}
}
