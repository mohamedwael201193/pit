package cli

import (
	"strings"
	"testing"

	"github.com/mohamedwael201193/pit/internal/identity"
)

func TestBindSameWalletIsStable(t *testing.T) {
	dir := t.TempDir()
	a := "0x1111111111111111111111111111111111111111"
	first, err := Bind(dir, "mainnet", a)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := identity.ParseWorkspaceID(first.WorkspaceID); err != nil {
		t.Fatal(err)
	}
	second, err := Bind(dir, "mainnet", strings.ToUpper(a))
	if err != nil {
		t.Fatal(err)
	}
	if first.WorkspaceID != second.WorkspaceID {
		t.Fatal("workspace must not rotate")
	}
}

func TestBindSecondWalletDenied(t *testing.T) {
	dir := t.TempDir()
	if _, err := Bind(dir, "mainnet", "0x1111111111111111111111111111111111111111"); err != nil {
		t.Fatal(err)
	}
	if _, err := Bind(dir, "mainnet", "0x2222222222222222222222222222222222222222"); err == nil {
		t.Fatal("user B")
	} else if err.Error() != "workspace_owned" {
		t.Fatal(err)
	}
}

func TestBindNetworkSwitchDenied(t *testing.T) {
	dir := t.TempDir()
	w := "0x1111111111111111111111111111111111111111"
	if _, err := Bind(dir, "mainnet", w); err != nil {
		t.Fatal(err)
	}
	if _, err := Bind(dir, "testnet", w); err == nil {
		t.Fatal("mixed network")
	}
}

func TestBindAfterForgetAllowsSecondWallet(t *testing.T) {
	dir := t.TempDir()
	if _, err := Bind(dir, "testnet", "0x1111111111111111111111111111111111111111"); err != nil {
		t.Fatal(err)
	}
	if err := Logout(dir, true); err != nil {
		t.Fatal(err)
	}
	st, err := Bind(dir, "testnet", "0x2222222222222222222222222222222222222222")
	if err != nil {
		t.Fatal(err)
	}
	if st.Wallet != "0x2222222222222222222222222222222222222222" {
		t.Fatal(st.Wallet)
	}
}

func TestPublicBindOmitsSecrets(t *testing.T) {
	body := PublicBind(DiskState{
		WorkspaceID: identity.NewWorkspaceID(),
		Network:     "mainnet",
		Wallet:      "0x1111111111111111111111111111111111111111",
	})
	if body["sign"] != false || body["trade"] != false {
		t.Fatal(body)
	}
}
