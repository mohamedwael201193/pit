package chain8004

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mohamedwael201193/pit/internal/config"
)

func TestProveFillRequiresOID(t *testing.T) {
	raw := json.RawMessage(`[{"oid":"531667200134","coin":"HYPE","sz":"0.16","px":"80.826","side":"B"}]`)
	row, err := ProveFill(raw, RecordedOID)
	if err != nil {
		t.Fatal(err)
	}
	if row.Coin != "HYPE" {
		t.Fatal(row)
	}
	if _, err := ProveFill(raw, "1"); err == nil {
		t.Fatal("missing oid")
	}
}

func TestCanonicalFeedbackHashStable(t *testing.T) {
	card, err := CanonicalFeedback("eip155:16661:0x8004A169FB4a3325136EB29fA0ceB6D2e539a432", MainnetAgentID, RecordedOID, "hype_fill", "successful_job")
	if err != nil {
		t.Fatal(err)
	}
	h1, b1, err := FeedbackHash(card)
	if err != nil {
		t.Fatal(err)
	}
	h2, _, err := FeedbackHash(card)
	if err != nil {
		t.Fatal(err)
	}
	if h1 != h2 {
		t.Fatal("unstable")
	}
	if strings.Contains(strings.ToLower(string(b1)), `"book"`) {
		t.Fatal("book")
	}
	if crypto.Keccak256Hash(b1) != h1 {
		t.Fatal("keccak")
	}
}

func TestCanonicalFeedbackUnknownJob(t *testing.T) {
	if _, err := CanonicalFeedback("eip155:16661:0x8004A169FB4a3325136EB29fA0ceB6D2e539a432", MainnetAgentID, "1", "hype_fill", "successful_job"); err == nil {
		t.Fatal("unknown")
	}
}

func TestRefuseSessionAgentReporter(t *testing.T) {
	if err := refuseSessionSigner(common.HexToAddress(RecordedHLAgent)); err == nil {
		t.Fatal("session")
	}
	if err := refuseSessionSigner(common.HexToAddress(RecordedOwner)); err != nil {
		t.Fatal(err)
	}
}

func TestLiveIdentityAndReputation(t *testing.T) {
	if os.Getenv("PIT_LIVE_IDENTITY") != "1" {
		t.Skip("set PIT_LIVE_IDENTITY=1")
	}
	ch := config.For(config.Mainnet)
	live, err := IdentitySnapshot(ch, MainnetAgentID)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.EqualFold(live.Owner, RecordedOwner) {
		t.Fatalf("owner %s", live.Owner)
	}
	reg, err := IdentityRegistry(ch)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.EqualFold(reg, ch.Identity8004) {
		t.Fatalf("registry %s", reg)
	}
	impl, err := Implementation(ch.RPC, ch.Reputation8004)
	if err != nil {
		t.Fatal(err)
	}
	if impl == "" || strings.EqualFold(impl, "0x0000000000000000000000000000000000000000") {
		t.Fatal("impl")
	}
}
