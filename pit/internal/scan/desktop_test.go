package scan

import (
	"path/filepath"
	"testing"
)

func TestDesktopSourceHasNoSessionMaterial(t *testing.T) {
	root := filepath.Join("..", "..", "..", "apps", "desktop", "src")
	if err := DesktopSource(root); err != nil {
		t.Fatal(err)
	}
}
