package companion

import (
	"strings"
	"testing"
	"time"

	"github.com/mohamedwael201193/pit/internal/auto"
)

func TestHuntSkipSetIncludesStoodDown(t *testing.T) {
	h := New(t.TempDir())
	h.job.coin = "AVAX"
	h.job.deny = "no_side"
	skip := h.huntSkipSet()
	if skip["AVAX"] == "" {
		t.Fatalf("expected AVAX skip, got %#v", skip)
	}
}

func TestHuntSkipSetIncludesPriorSlice(t *testing.T) {
	h := New(t.TempDir())
	h.huntSkip = []string{"AVAX", "SOL"}
	skip := h.huntSkipSet()
	if skip["AVAX"] == "" || skip["SOL"] == "" {
		t.Fatalf("expected prior hunt skips, got %#v", skip)
	}
}

func TestResolveChatCoinDoesNotWrap(t *testing.T) {
	skip := map[string]string{"AVAX": "stood_down", "HYPE": "stood_down"}
	if got := resolveChatCoin("AVAX", skip, "", false); got != "" {
		t.Fatalf("wrap %q", got)
	}
	if got := resolveChatCoin("AVAX", skip, "DOGE", false); got != "DOGE" {
		t.Fatalf("next %q", got)
	}
	if got := resolveChatCoin("AVAX", skip, "", true); got != "AVAX" {
		t.Fatalf("fresh %q", got)
	}
	if got := resolveChatCoin("", skip, "SOL", false); got != "SOL" {
		t.Fatalf("unnamed %q", got)
	}
}

func TestThisHuntSkipIgnoresStaleAuto(t *testing.T) {
	dir := t.TempDir()
	p := auto.Load(dir)
	p.RememberSkip("BTC", "no_side", auto.SkipStoodDown, 4*time.Hour)
	p.RememberSkip("ETH", "no_side", auto.SkipStoodDown, 4*time.Hour)
	p.RememberSkip("SOL", "no_side", auto.SkipStoodDown, 4*time.Hour)
	p.RememberSkip("HYPE", "challenger_killed", auto.SkipStoodDown, 4*time.Hour)
	if err := auto.Save(dir, p); err != nil {
		t.Fatal(err)
	}
	h := New(dir)
	h.huntSkip = []string{"DOGE", "AVAX"}
	this := h.thisHuntSkipSet()
	if this["BTC"] != "" || this["ETH"] != "" || this["SOL"] != "" || this["HYPE"] != "" {
		t.Fatalf("chat hunt must not inherit 4h auto skips, got %#v", this)
	}
	if this["DOGE"] == "" || this["AVAX"] == "" {
		t.Fatalf("this hunt skips missing, got %#v", this)
	}
	if got := resolveChatCoin("BTC", this, "ETH", false); got != "BTC" {
		t.Fatalf("named remaining book %q", got)
	}
	merged := h.huntSkipSet()
	if merged["BTC"] == "" {
		t.Fatal("automation skip set still includes 4h BTC")
	}
}

func TestChatSourceDoesNotAutoContinue(t *testing.T) {
	if normalizeResearchSource("chat") == "automation" {
		t.Fatal("chat must not kick autoTick")
	}
	if normalizeResearchSource("research_ui") == "automation" {
		t.Fatal("research_ui must not kick autoTick")
	}
	if normalizeResearchSource("automation") != "automation" {
		t.Fatal("automation keeps autoTick")
	}
}

func TestHuntSkipPersistsAndMergesAuto(t *testing.T) {
	dir := t.TempDir()
	p := auto.Load(dir)
	p.RememberSkip("AVAX", "no_side", auto.SkipStoodDown, time.Hour)
	if err := auto.Save(dir, p); err != nil {
		t.Fatal(err)
	}
	h := New(dir)
	h.researchMu.Lock()
	h.huntSkip = []string{"HYPE"}
	h.persistHuntSkipLocked()
	h.researchMu.Unlock()
	h2 := New(dir)
	skip := h2.huntSkipSet()
	if skip["AVAX"] == "" || skip["HYPE"] == "" {
		t.Fatalf("expected persisted + auto skip, got %#v", skip)
	}
}

func TestNamedRolesVerifiedRequiresThree(t *testing.T) {
	one := []map[string]any{{"role": "researcher", "verify_e2ee": "OK"}}
	if namedRolesVerified(one) {
		t.Fatal("one role is not a committee")
	}
	if TerminalKind(false, "", "", false, false, one) != TermCommitteeIncomplete {
		t.Fatalf("got %q", TerminalKind(false, "", "", false, false, one))
	}
	three := []map[string]any{
		{"role": "researcher", "verify_e2ee": "OK"},
		{"role": "challenger", "verify_e2ee": "OK"},
		{"role": "risk", "verify_e2ee": "OK"},
	}
	if !namedRolesVerified(three) {
		t.Fatal("three OK")
	}
	if TerminalKind(false, "", "", true, true, three) != TermReadyEligible {
		t.Fatal(TerminalKind(false, "", "", true, true, three))
	}
}

func TestReplyWhyNotContinuesPastOneBook(t *testing.T) {
	h := &Hub{Dir: t.TempDir()}
	got := h.replyWhyNotTrade()
	if strings.Contains(got, "Latest: BTC") {
		t.Fatal(got)
	}
	if !strings.Contains(got, "Scan continues") {
		t.Fatal(got)
	}
	if strings.Contains(got, `"execute":true`) {
		t.Fatal(got)
	}
}

func TestGuardedAlreadyAttemptedSameHash(t *testing.T) {
	if guardedAlreadyAttempted(nil, "0xabc") {
		t.Fatal("nil")
	}
	if !guardedAlreadyAttempted(map[string]any{"hash": "0xabc", "oid": "", "posted": false}, "0xabc") {
		t.Fatal("same hash with missing oid must not retry POST")
	}
	if guardedAlreadyAttempted(map[string]any{"hash": "0xdef"}, "0xabc") {
		t.Fatal("other hash")
	}
}

func TestReplyResearchStandDownIsSuccessfulNoTrade(t *testing.T) {
	h := &Hub{Dir: t.TempDir()}
	h.job.deny = "no_side"
	h.job.eligible = false
	h.job.roles = []map[string]any{
		{"role": "researcher", "verify_e2ee": "OK"},
		{"role": "challenger", "verify_e2ee": "OK"},
		{"role": "risk", "verify_e2ee": "OK"},
	}
	got := h.replyResearch("")
	if !strings.Contains(got, "verified no-trade") || !strings.Contains(got, "next eligible") {
		t.Fatal(got)
	}
	if strings.Contains(got, `"execute":true`) {
		t.Fatal(got)
	}
}

func TestReplyResearchBelowMinIsNotStandDown(t *testing.T) {
	h := &Hub{Dir: t.TempDir()}
	h.job.deny = "below_min_notional"
	h.job.eligible = false
	got := h.replyResearch("")
	if strings.Contains(got, "Committee stood down") {
		t.Fatal(got)
	}
	if !strings.Contains(got, "not sizeable") {
		t.Fatal(got)
	}
}

func TestStandDownIsNotCrash(t *testing.T) {
	roles := []map[string]any{
		{"role": "researcher", "verify_e2ee": "OK"},
		{"role": "challenger", "verify_e2ee": "OK"},
		{"role": "risk", "verify_e2ee": "OK"},
	}
	if TerminalKind(false, "risk_killed", "risk_killed", true, false, roles) != TermReadyStoodDown {
		t.Fatal(TerminalKind(false, "risk_killed", "risk_killed", true, false, roles))
	}
	if researchCardTitle(TermReadyStoodDown) == "RESEARCH COMPLETE" {
		t.Fatal("title")
	}
}

func TestTimeoutIsNotTEE(t *testing.T) {
	if TerminalKind(false, "DIRECT_PROVIDER_TIMEOUT", "", false, false, nil) != TermDirectTimeout {
		t.Fatal(TerminalKind(false, "DIRECT_PROVIDER_TIMEOUT", "", false, false, nil))
	}
	if TerminalKind(false, "TEE_VERIFY_FAIL", "", false, false, nil) != TermTEESignatureInvalid {
		t.Fatal("tee")
	}
}

func TestPollFailedNeverTerminatesRunningJob(t *testing.T) {
	if TerminalKind(true, "POLL_FAILED", "", false, false, nil) != "" {
		t.Fatal("poll while running")
	}
}

func TestCanceledAndPolicyAndMarket(t *testing.T) {
	if TerminalKind(false, "research_cancelled", "", false, false, nil) != TermCanceledByUser {
		t.Fatal("cancel")
	}
	if TerminalKind(false, "asset_not_allowed", "", false, false, nil) != TermPolicyDenied {
		t.Fatal("policy")
	}
	if TerminalKind(false, "empty_envelope", "", false, false, nil) != TermMarketDenied {
		t.Fatal("market")
	}
	if TerminalKind(false, "JOB_CRASHED", "", false, false, nil) != TermJobCrashed {
		t.Fatal("crash")
	}
}

func TestRateLimitAndSponsorAndSigner(t *testing.T) {
	if TerminalKind(false, "429", "", false, false, nil) != TermDirectRateLimited {
		t.Fatal("429")
	}
	if TerminalKind(false, TermSponsorQuota, "", false, false, nil) != TermSponsorQuota {
		t.Fatal("quota")
	}
	if TerminalKind(false, TermWrongNetwork, "", false, false, nil) != TermWrongNetwork {
		t.Fatal("net")
	}
	if TerminalKind(false, TermTEESignerMismatch, "", false, false, nil) != TermTEESignerMismatch {
		t.Fatal("signer")
	}
	if TerminalKind(false, TermDirectUnavailable, "", false, false, nil) != TermDirectUnavailable {
		t.Fatal("unavail")
	}
}

func TestVerifiedTitleOnlyForEligibleCommittee(t *testing.T) {
	if researchCardTitle(TermReadyEligible) != "RESEARCH COMPLETE" {
		t.Fatal("title")
	}
	if researchCardTitle(TermCommitteeIncomplete) == "RESEARCH COMPLETE" {
		t.Fatal("incomplete")
	}
}
