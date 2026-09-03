package chainrpc

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
)

func TestAddressOfKeyMatchesPubkey(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	got, err := AddressOfKey(hex.EncodeToString(crypto.FromECDSA(key)))
	if err != nil {
		t.Fatal(err)
	}
	want := crypto.PubkeyToAddress(key.PublicKey)
	if got != want {
		t.Fatalf("%s != %s", got.Hex(), want.Hex())
	}
}

func TestSignLegacyRoundTrip(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	to := crypto.PubkeyToAddress(key.PublicKey)
	raw, err := signLegacyEIP155(key, 16661, 0, big.NewInt(1), 21000, to, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(raw) < 32 {
		t.Fatal("short")
	}
}
