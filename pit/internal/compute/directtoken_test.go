package compute

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
)

func TestCanonicalJSONMatchesJSFieldOrder(t *testing.T) {
	tok := SessionToken{
		Address:    "0x1111111111111111111111111111111111111111",
		Provider:   "0x7DCFe6AEa70350C2090041524c9B4A9262DCe87D",
		Timestamp:  1,
		ExpiresAt:  2,
		Nonce:      "aabbccdd",
		Generation: 0,
		TokenId:    255,
	}
	got := CanonicalJSON(tok)
	want := `{"address":"0x1111111111111111111111111111111111111111","provider":"0x7DCFe6AEa70350C2090041524c9B4A9262DCe87D","timestamp":1,"expiresAt":2,"nonce":"aabbccdd","generation":0,"tokenId":255}`
	if got != want {
		t.Fatalf("got %s", got)
	}
	var round SessionToken
	if err := json.Unmarshal([]byte(got), &round); err != nil {
		t.Fatal(err)
	}
	if CanonicalJSON(round) != got {
		t.Fatal("roundtrip")
	}
}

func TestAssembleAndRecoverOfficialPrefix(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	wallet := crypto.PubkeyToAddress(key.PublicKey).Hex()
	_, ch, err := NewChallenge(wallet, "0x7DCFe6AEa70350C2090041524c9B4A9262DCe87D", 0, time.UnixMilli(1_700_000_000_000))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(ch.Digest, "0x") || len(ch.Digest) != 66 {
		t.Fatal(ch.Digest)
	}
	messageHash := crypto.Keccak256Hash([]byte(ch.Message))
	prefixed := crypto.Keccak256Hash([]byte("\x19Ethereum Signed Message:\n32"), messageHash.Bytes())
	sig, err := crypto.Sign(prefixed.Bytes(), key)
	if err != nil {
		t.Fatal(err)
	}
	sig[64] += 27
	auth, tok, err := AcceptDirectSignature(ch.Message, "0x"+encodeHex(sig), wallet, time.UnixMilli(1_700_000_000_000))
	if err != nil {
		t.Fatal(err)
	}
	if err := RefuseRouterKey(auth); err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(auth, "Bearer app-sk-") {
		t.Fatal(auth)
	}
	if strings.Contains(strings.ToLower(ch.Explain), "seed") && strings.Contains(ch.Explain, "private key") {
		t.Fatal("must not ask for a key")
	}
	parsed, _, err := ParseBearer(auth)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.TokenId != 255 || parsed.Address == "" {
		t.Fatal(parsed)
	}
	if TokenExpired(tok, time.UnixMilli(1_700_000_000_000)) {
		t.Fatal("fresh")
	}
	if !TokenExpired(tok, time.UnixMilli(tok.ExpiresAt+1000)) {
		t.Fatal("expired")
	}
}

func TestAcceptDirectSignatureWrongWallet(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	wallet := crypto.PubkeyToAddress(key.PublicKey).Hex()
	_, ch, err := NewChallenge(wallet, "0x7DCFe6AEa70350C2090041524c9B4A9262DCe87D", 0, time.UnixMilli(1_700_000_000_000))
	if err != nil {
		t.Fatal(err)
	}
	messageHash := crypto.Keccak256Hash([]byte(ch.Message))
	prefixed := crypto.Keccak256Hash([]byte("\x19Ethereum Signed Message:\n32"), messageHash.Bytes())
	sig, err := crypto.Sign(prefixed.Bytes(), key)
	if err != nil {
		t.Fatal(err)
	}
	sig[64] += 27
	_, _, err = AcceptDirectSignature(ch.Message, "0x"+encodeHex(sig), "0x2222222222222222222222222222222222222222", time.UnixMilli(1_700_000_000_000))
	if err == nil || err.Error() != "signature_mismatch" {
		t.Fatalf("%v", err)
	}
}

func TestAssembleRefuseRouter(t *testing.T) {
	if err := RefuseRouterKey("Bearer sk-abc"); err == nil {
		t.Fatal("sk")
	}
}

func encodeHex(b []byte) string {
	const hexdigits = "0123456789abcdef"
	out := make([]byte, len(b)*2)
	for i, v := range b {
		out[i*2] = hexdigits[v>>4]
		out[i*2+1] = hexdigits[v&0x0f]
	}
	return string(out)
}
