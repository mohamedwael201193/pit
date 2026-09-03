package chain8004

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mohamedwael201193/pit/internal/config"
	"github.com/mohamedwael201193/pit/internal/identity"
)

func TestSubmitFillFeedbackRefusesOwner(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	owner := crypto.PubkeyToAddress(key.PublicKey)
	own, _ := identity.NormalizeAddress(owner.Hex())
	if err := FeedbackAllowed(own, own, own); err == nil {
		t.Fatal("self")
	}
	ch := config.For(config.Mainnet)
	_, err = SubmitFillFeedback(ch, "0x"+common.Bytes2Hex(crypto.FromECDSA(key)), MainnetAgentID, RecordedOID, []byte(`[]`))
	if err == nil {
		t.Fatal("empty fills")
	}
}

func TestSetURIRefusesBadKey(t *testing.T) {
	ch := config.For(config.Mainnet)
	if _, err := SetURIFromOwner(ch, "0xab", MainnetAgentID, AgentCardURL); err == nil {
		t.Fatal("key")
	}
}
