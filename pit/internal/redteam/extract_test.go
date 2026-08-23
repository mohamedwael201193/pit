package redteam

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBrowserExtract(t *testing.T) {
	dir := t.TempDir()
	ok := filepath.Join(dir, "ok.tsx")
	if err := os.WriteFile(ok, []byte("export const copy = 'Connect your wallet'"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := BrowserExtract(dir); err != nil {
		t.Fatal(err)
	}
	bad := filepath.Join(dir, "bad.tsx")
	if err := os.WriteFile(bad, []byte("const session_key = 'x'"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := BrowserExtract(dir); err == nil {
		t.Fatal("extract")
	}
}
