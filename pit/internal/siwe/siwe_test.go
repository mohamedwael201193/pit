package siwe

import (
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mohamedwael201193/pit/internal/identity"
)

func encode(b []byte) string {
	const hex = "0123456789abcdef"
	out := make([]byte, len(b)*2)
	for i, v := range b {
		out[i*2] = hex[v>>4]
		out[i*2+1] = hex[v&0x0f]
	}
	return string(out)
}

func TestRecoverRoundTrip(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	addr := crypto.PubkeyToAddress(key.PublicKey)
	want, err := identity.NormalizeAddress(addr.Hex())
	if err != nil {
		t.Fatal(err)
	}
	msg, err := Build(Message{
		Domain:   "localhost:3000",
		Address:  want,
		URI:      "http://localhost:3000",
		ChainID:  16661,
		Nonce:    "pitnonce1",
		IssuedAt: time.Now(),
	})
	if err != nil {
		t.Fatal(err)
	}
	h := accountsHash([]byte(msg))
	sig, err := crypto.Sign(h, key)
	if err != nil {
		t.Fatal(err)
	}
	sig[64] += 27
	got, err := Recover(msg, "0x"+encode(sig))
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("got %s want %s", got, want)
	}
}

func TestBadSignature(t *testing.T) {
	if _, err := Recover("hello", "0x00"); err == nil {
		t.Fatal("expected SIGNATURE_DECLINED")
	}
}

func TestWrongNetwork(t *testing.T) {
	if err := AssertChain(1, 16661); err == nil || err.Error() != "WRONG_NETWORK" {
		t.Fatalf("got %v", err)
	}
}
