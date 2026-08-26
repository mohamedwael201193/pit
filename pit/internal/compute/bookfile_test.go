package compute

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadEnvelopeRequiresBothFiles(t *testing.T) {
	if _, _, err := LoadEnvelope("", ""); err == nil {
		t.Fatal("empty")
	}
	dir := t.TempDir()
	m := filepath.Join(dir, "m.json")
	b := filepath.Join(dir, "b.json")
	if err := os.WriteFile(m, []byte(`{"venue":"hyperliquid","coin":"ETH"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(b, []byte(`{"positions":[]}`), 0o600); err != nil {
		t.Fatal(err)
	}
	market, book, err := LoadEnvelope(m, b)
	if err != nil {
		t.Fatal(err)
	}
	if len(market) == 0 || len(book) == 0 {
		t.Fatal("bytes")
	}
}

func TestLoadEnvelopeBlankFile(t *testing.T) {
	dir := t.TempDir()
	m := filepath.Join(dir, "m.json")
	b := filepath.Join(dir, "b.json")
	if err := os.WriteFile(m, []byte("   "), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(b, []byte(`{"positions":[]}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, _, err := LoadEnvelope(m, b); err == nil {
		t.Fatal("blank")
	}
}
