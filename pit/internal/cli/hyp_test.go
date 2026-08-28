package cli

import (
	"os"
	"testing"
)

func TestSaveLoadHypothesis(t *testing.T) {
	dir := t.TempDir()
	if LoadHypothesis(dir) != "none" {
		t.Fatal("default")
	}
	if err := SaveHypothesis(dir, "long"); err != nil {
		t.Fatal(err)
	}
	if LoadHypothesis(dir) != "long" {
		t.Fatal(LoadHypothesis(dir))
	}
	if err := SaveHypothesis(dir, "withdraw"); err == nil {
		t.Fatal("withdraw")
	}
	raw, _ := os.ReadFile(hypothesisFile(dir))
	if string(raw) == "" {
		t.Fatal("missing")
	}
}
