package cli

import (
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mohamedwael201193/pit/internal/compute"
	"github.com/mohamedwael201193/pit/internal/identity"
)

func TestDirectChallengeRoundTripKeychain(t *testing.T) {
	t.Setenv("PIT_KEYRING", "file")
	t.Setenv("PIT_DIRECT_AUTH_FILE", "")
	t.Setenv("PIT_DIRECT_SPONSOR_FILE", "")
	dir := t.TempDir()
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	wallet := crypto.PubkeyToAddress(key.PublicKey).Hex()
	ws := identity.NewWorkspaceID()
	if err := Save(dir, DiskState{WorkspaceID: ws, Network: "mainnet", Wallet: strings.ToLower(wallet)}); err != nil {
		t.Fatal(err)
	}
	ch, err := IssueDirectChallenge(dir)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(strings.ToLower(ch.Message), "app-sk-") {
		t.Fatal("challenge leak")
	}
	messageHash := crypto.Keccak256Hash([]byte(ch.Message))
	prefixed := crypto.Keccak256Hash([]byte("\x19Ethereum Signed Message:\n32"), messageHash.Bytes())
	sig, err := crypto.Sign(prefixed.Bytes(), key)
	if err != nil {
		t.Fatal(err)
	}
	sig[64] += 27
	meta, err := CompleteDirect(dir, "", "0x"+encodeSig(sig))
	if err != nil {
		t.Fatal(err)
	}
	if meta.Source != "keychain" || meta.TokenId != 255 {
		t.Fatal(meta)
	}
	file, got, err := ResolveWorkspaceAuth(dir)
	if err != nil {
		t.Fatal(err)
	}
	if err := compute.RefuseRouterKey(file.Authorization); err != nil {
		t.Fatal(err)
	}
	if got.Source != "keychain" {
		t.Fatal(got)
	}
	st := DirectStatus(dir)
	if st["ok"] != true || strings.Contains(strings.ToLower(st["source"].(string)), "app-sk") {
		t.Fatal(st)
	}
	c := checkDirectAuth(dir)
	if !c.OK {
		t.Fatal(c)
	}
	ForgetDirect(dir)
	if _, _, err := ResolveWorkspaceAuth(dir); err == nil {
		t.Fatal("deleted")
	}
	_ = time.Now()
}

func TestDirectWrongWalletRefused(t *testing.T) {
	t.Setenv("PIT_KEYRING", "file")
	dir := t.TempDir()
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	other, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	wallet := crypto.PubkeyToAddress(key.PublicKey).Hex()
	ws := identity.NewWorkspaceID()
	if err := Save(dir, DiskState{WorkspaceID: ws, Network: "mainnet", Wallet: strings.ToLower(wallet)}); err != nil {
		t.Fatal(err)
	}
	ch, err := IssueDirectChallenge(dir)
	if err != nil {
		t.Fatal(err)
	}
	messageHash := crypto.Keccak256Hash([]byte(ch.Message))
	prefixed := crypto.Keccak256Hash([]byte("\x19Ethereum Signed Message:\n32"), messageHash.Bytes())
	sig, err := crypto.Sign(prefixed.Bytes(), other)
	if err != nil {
		t.Fatal(err)
	}
	sig[64] += 27
	if _, err := CompleteDirect(dir, "", "0x"+encodeSig(sig)); err == nil {
		t.Fatal("wrong wallet")
	}
}

func encodeSig(b []byte) string {
	const d = "0123456789abcdef"
	out := make([]byte, len(b)*2)
	for i, v := range b {
		out[i*2] = d[v>>4]
		out[i*2+1] = d[v&0x0f]
	}
	return string(out)
}
