package scan

import (
	"path/filepath"
	"testing"
)

func TestWebSourceHasNoSessionMaterial(t *testing.T) {
	root := filepath.Join("..", "..", "..", "apps", "web", "src")
	if err := WebSource(root); err != nil {
		t.Fatal(err)
	}
}
