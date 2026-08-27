package hl

import (
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
)

func TestSignL1RecoversAgent(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	raw, err := BuildOrder(1, true, "2500", "0.004", "0x11111111111111111111111111111111")
	if err != nil {
		t.Fatal(err)
	}
	env, err := SignL1(key, raw, 1700000000000, false)
	if err != nil {
		t.Fatal(err)
	}
	if !env.Signed() {
		t.Fatal("unsigned")
	}
	got, err := RecoverL1(env, false)
	if err != nil {
		t.Fatal(err)
	}
	want := strings.ToLower(crypto.PubkeyToAddress(key.PublicKey).Hex())
	if got != want {
		t.Fatalf("%s %s", got, want)
	}
}
