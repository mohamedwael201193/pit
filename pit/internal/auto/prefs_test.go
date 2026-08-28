package auto

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRefuseExecute(t *testing.T) {
	if RefuseExecute() == nil {
		t.Fatal("must refuse")
	}
}

func TestSaveStripsExecute(t *testing.T) {
	dir := t.TempDir()
	p := Default()
	p.Execute = true
	p.AutoResearch = true
	if err := Save(dir, p); err != nil {
		t.Fatal(err)
	}
	got := Load(dir)
	if got.Execute {
		t.Fatal("execute leaked")
	}
	if !got.AutoResearch || !got.Watch {
		t.Fatalf("%+v", got)
	}
	raw, _ := os.ReadFile(filepath.Join(dir, "automation.json"))
	if string(raw) == "" {
		t.Fatal("missing")
	}
}

func TestMatchesTrigger(t *testing.T) {
	if !Matches("policy_pass", "in_universe") {
		t.Fatal("universe")
	}
	if Matches("funding", "mark_below_oracle") {
		t.Fatal("funding mismatch")
	}
	if !Matches("mark_below_oracle", "mark_below_oracle") {
		t.Fatal("gap")
	}
}
