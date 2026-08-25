package mcp

import "testing"

func TestAuthFileNever(t *testing.T) {
	if !AuthFileNever() {
		t.Fatal("auth_file")
	}
}

func TestSealerNever(t *testing.T) {
	if !SealerNever() {
		t.Fatal("sealer")
	}
}
