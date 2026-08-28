package cli

import (
	"strings"
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
	if !DoctorFailed([]Check{w, {Name: "network"}, {Name: "keychain"}, {Name: "hyperliquid"}, {Name: "0g_rpc"}, {Name: "memory_key"}}) {
		t.Fatal("required")
	}
}

func TestDoctorRefusesGlobalMemoryKey(t *testing.T) {
	t.Setenv("PIT_MEMORY_KEY", "0xabc")
	c := checkMemoryEnv()
	if c.OK {
		t.Fatal("must refuse global memory key")
	}
}

func TestDoctorDirectAuthUnset(t *testing.T) {
	t.Setenv("PIT_KEYRING", "file")
	t.Setenv("PIT_DIRECT_AUTH_FILE", "")
	t.Setenv("PIT_DIRECT_SPONSOR_FILE", "")
	c := checkDirectAuth(t.TempDir())
	if c.OK {
		t.Fatal("unset auth must not pass")
	}
	if c.Detail == "PIT_DIRECT_AUTH_FILE unset" {
		t.Fatal("must not send a normal user to an env file")
	}
}

func TestDoctorVersionPresent(t *testing.T) {
	c := checkVersion()
	if !c.OK || c.Detail == "" {
		t.Fatal(c)
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

func TestDoctorHLAgentUnbound(t *testing.T) {
	c := checkHLAgent(t.TempDir())
	if c.OK {
		t.Fatal("unbound agent must wait")
	}
	if c.Name != "hl_agent" {
		t.Fatal(c.Name)
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

func TestDoctorTeeMissingIsHonest(t *testing.T) {
	c := checkTee(t.TempDir())
	if c.OK {
		t.Fatal("empty last-research must wait")
	}
	if c.Name != "tee" {
		t.Fatal(c.Name)
	}
}

func TestDoctorCreditUnbound(t *testing.T) {
	c := checkDirectCredit(t.TempDir())
	if c.OK {
		t.Fatal("unbound credit")
	}
}

func TestDoctorCreditUnreadNotZero(t *testing.T) {
	dir := t.TempDir()
	id := identity.NewWorkspaceID()
	if err := Save(dir, DiskState{WorkspaceID: id, Network: "mainnet", Wallet: "0x1111111111111111111111111111111111111111"}); err != nil {
		t.Fatal(err)
	}
	c := checkDirectCredit(dir)
	if strings.Contains(strings.ToLower(c.Detail), "unread") && strings.Contains(c.Detail, "0 0G") {
		t.Fatal(c.Detail)
	}
}

func TestCompanionDoctorDoesNotSelfDial(t *testing.T) {
	t.Setenv("PIT_COMPANION", "1")
	c := checkCompanion()
	if !c.OK || c.Detail != "this process" {
		t.Fatal(c)
	}
}
