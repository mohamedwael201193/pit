package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadLastOrderKeepsLargeOID(t *testing.T) {
	dir := t.TempDir()
	raw := []byte(`{"oid":529167222216,"status":"filled","cloid":"0xacb8a6b8a476fbb7cbeccf18d78b058c","sign":false,"trade":false}`)
	if err := os.WriteFile(filepath.Join(dir, "last-order.json"), raw, 0o600); err != nil {
		t.Fatal(err)
	}
	got := LoadLastOrder(dir)
	if got["oid"] != "529167222216" {
		t.Fatalf("oid %v", got["oid"])
	}
	if got["status"] != "filled" {
		t.Fatalf("status %v", got["status"])
	}
}
