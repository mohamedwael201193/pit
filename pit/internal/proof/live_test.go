package proof

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/mohamedwael201193/pit/internal/config"
	"github.com/mohamedwael201193/pit/internal/receipt"
	"github.com/mohamedwael201193/pit/internal/storage"
)

// TestLiveFileAndVerify publishes a real receipt to 0G Storage and re-checks it
// the way a third party would. It is skipped unless the operator opts in because
// it spends 0G gas.
func TestLiveFileAndVerify(t *testing.T) {
	if os.Getenv("PIT_LIVE_PROOF") != "1" {
		t.Skip("set PIT_LIVE_PROOF=1 to file real evidence on 0G")
	}
	ch := config.For(config.Mainnet)
	f := Filer{
		CLI: storage.LookCLI(), RPC: ch.RPC, Indexer: ch.StorageIndexer, Flow: ch.StorageFlow,
		Explorer: ch.Explorer, Network: string(ch.Network), ChainID: ch.ChainID,
		PayerKey: strings.TrimSpace(os.Getenv("PIT_OG_PAYER_KEY")), Dir: t.TempDir(),
	}
	if err := f.Ready(); err != nil {
		t.Fatal(err)
	}
	r := receipt.New(receipt.KindResearch, f.Network, f.ChainID, "ws-live-proof-check")
	r.Venue = "hyperliquid"
	r.Market = "ETH"
	r.Verdict = "READY_STOOD_DOWN"
	r.Deny = "challenger_killed"
	r.PreviewHash = "0x" + strings.Repeat("ab", 32)
	r.PolicyHash = "0x" + strings.Repeat("cd", 32)
	r.Roles = []receipt.Role{
		{Role: "researcher", VerifyE2EE: "OK", Signer: "0xabc", TeeSigner: "0xabc", Side: "long"},
		{Role: "challenger", VerifyE2EE: "OK", Signer: "0xabc", TeeSigner: "0xabc", Kill: true},
		{Role: "risk", VerifyE2EE: "OK", Signer: "0xabc", TeeSigner: "0xabc", Survives: true},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Minute)
	defer cancel()
	filed, err := f.File(ctx, r)
	if err != nil {
		t.Fatal(err)
	}
	if err := storage.RejectBadRoot(filed.Root); err != nil {
		t.Fatal("no usable root")
	}
	t.Log("root", filed.Root)
	t.Log("digest", filed.Digest)
	t.Log("tx", filed.Tx)
	t.Log("tx link", filed.TxLink)
	if filed.Tx == "" && !filed.Duplicate {
		t.Fatal("a fresh filing must record a chain transaction")
	}
	rows, err := Index(f.Dir, 0)
	if err != nil || len(rows) != 1 {
		t.Fatalf("index did not persist: %v rows=%d", err, len(rows))
	}
	got, err := f.Verify(ctx, filed.Root)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("verify %+v", got)
	if got.Failure != "" {
		t.Fatal("verification failed:", got.Failure)
	}
	if !got.ProofValidated {
		t.Fatal("merkle proof was not validated")
	}
	if !got.DigestMatch || !strings.EqualFold(got.Recomputed, filed.Digest) {
		t.Fatal("digest did not survive the round trip")
	}
	if !got.PublicSafe || !got.RolesVerified {
		t.Fatal("published record failed its own honesty checks")
	}
	if got.FinalizedNodes == 0 {
		t.Fatal("no storage node reported a finalized copy")
	}
	if filed.Tx != "" {
		if !got.AnchorBound || got.Anchor == nil {
			t.Fatal("chain transaction did not commit the root")
		}
		t.Log("anchor block", got.Anchor.BlockNumber, "flow", got.Anchor.Flow)
	}
	tampered, err := f.Verify(ctx, "0x"+strings.Repeat("ef", 32))
	if err == nil && tampered.Failure == "" {
		t.Fatal("a root that was never filed must not verify")
	}
}
