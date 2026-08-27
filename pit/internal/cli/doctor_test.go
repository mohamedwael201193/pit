package cli

import (
	"testing"

	"github.com/mohamedwael201193/pit/internal/identity"
)

func TestDoctorUnboundIsHonest(t *testing.T) {
	dir := t.TempDir()
	w := checkWallet(dir)
	s := checkSession(dir)
	if w.OK || s.OK {
		t.Fatal("unbound")
	}
	if !DoctorFailed([]Check{w, {Name: "network"}, {Name: "keychain"}, {Name: "hyperliquid"}, {Name: "0g_rpc"}}) {
		t.Fatal("required")
	}
}

func TestDoctorKeychainFile(t *testing.T) {
	c := checkKeychain(t.TempDir())
	if !c.OK {
		t.Fatal(c)
	}
	if c.Detail != "file" {
		t.Fatal(c.Detail)
	}
}

func TestWantJSONStripsFlag(t *testing.T) {
	want, rest := WantJSON([]string{"status", "--json", "--i-understand"})
	if !want || rest[0] != "status" || rest[1] != "--i-understand" {
		t.Fatalf("%v %v", want, rest)
	}
}

func TestDoctorWalletBound(t *testing.T) {
	dir := t.TempDir()
	id := identity.NewWorkspaceID()
	if err := Save(dir, DiskState{WorkspaceID: id, Network: "mainnet", Wallet: "0x1111111111111111111111111111111111111111"}); err != nil {
		t.Fatal(err)
	}
	c := checkWallet(dir)
	if !c.OK {
		t.Fatal(c)
	}
}
