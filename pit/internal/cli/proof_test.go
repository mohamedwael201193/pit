package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseProofFlagsRequireAll(t *testing.T) {
	if _, err := ParseProofFlags(nil); err == nil {
		t.Fatal("empty")
	}
	f, err := ParseProofFlags([]string{"--root", "0xabc1234567", "--out", "out.bin", "--key-file", "k"})
	if err != nil || f.Root != "0xabc1234567" {
		t.Fatal(err, f)
	}
}

func TestLoadProofKeyIgnoresEnv(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "k")
	if err := os.WriteFile(p, []byte("0x"+ "11"+string(make([]byte, 0))), 0o600); err != nil {
		t.Fatal(err)
	}
	// 32-byte hex
	hex := "0x" + "aa" + "bb" + "cc" + "dd" + "ee" + "ff" + "00" + "11" + "22" + "33" + "44" + "55" + "66" + "77" + "88" + "99"
	hex = "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	if err := os.WriteFile(p, []byte(hex+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	got, err := LoadProofKey(p, "0xdead")
	if err != nil || got != hex {
		t.Fatal(err, got)
	}
	if _, err := LoadProofKey("", "0xdead"); err == nil {
		t.Fatal("env key must be denied")
	}
}
