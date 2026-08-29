package demo

import (
	"os"
	"testing"
)

func TestReplayIsNeverLive(t *testing.T) {
	v := Replay(t.TempDir(), []map[string]any{{"root": "0xabc"}}, nil)
	if v.Live || v.Mode != "replay" || v.Sign || v.Trade {
		t.Fatalf("%+v", v)
	}
	if v.Count != 1 || v.Label == "" {
		t.Fatal(v)
	}
}

func TestLiveIsNotReplay(t *testing.T) {
	v := Live(t.TempDir())
	if !v.Live || v.Mode != "live" || v.Sign || v.Trade {
		t.Fatalf("%+v", v)
	}
}

func TestPrefRoundTrip(t *testing.T) {
	dir := t.TempDir()
	if LoadPref(dir) != "live" {
		t.Fatal("default")
	}
	if err := SavePref(dir, "replay"); err != nil {
		t.Fatal(err)
	}
	if LoadPref(dir) != "replay" {
		t.Fatal(LoadPref(dir))
	}
	raw, _ := os.ReadFile(PrefPath(dir))
	if string(raw) == "" {
		t.Fatal("empty")
	}
}
