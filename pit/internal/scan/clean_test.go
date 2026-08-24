package scan

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWebMustStayClean(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "App.tsx")
	if err := os.WriteFile(p, []byte("export function App(){return null}"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := WebMustStayClean(dir); err != nil {
		t.Fatal(err)
	}
}
