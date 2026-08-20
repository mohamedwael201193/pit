package session

import (
	"testing"

	"github.com/mohamedwael201193/pit/internal/keyring"
)

func TestKeygenNeverExports(t *testing.T) {
	k, exp, err := GenerateAgent(1)
	if err != nil {
		t.Fatal(err)
	}
	if k.Address == "" || exp == 0 {
		t.Fatal("empty")
	}
	if err := k.ExportJSON(); err == nil {
		t.Fatal("export must deny")
	}
	card := Card()
	if !card.Order || !card.Cancel || card.Withdraw || card.Leverage {
		t.Fatalf("%+v", card)
	}
	ring, err := keyring.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if err := k.Store(ring, "ws-a", "s1"); err != nil {
		t.Fatal(err)
	}
	got, err := LoadAgent(ring, "ws-a", "s1")
	if err != nil || got.Address != k.Address {
		t.Fatal(err)
	}
	if _, err := LoadAgent(ring, "ws-b", "s1"); err == nil {
		t.Fatal("cross workspace")
	}
	if k.Name("abcdef12-xxxx") != "PIT-abcdef12" && k.Name("abcdef12xxxx") != "PIT-abcdef1" {
		if !startsPIT(k.Name("workspace-id-here")) {
			t.Fatal(k.Name("workspace-id-here"))
		}
	}
}

func startsPIT(s string) bool { return len(s) > 4 && s[:4] == "PIT-" }
