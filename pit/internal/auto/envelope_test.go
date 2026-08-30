package auto

import (
	"strings"
	"testing"

	"github.com/mohamedwael201193/pit/internal/policy"
)

func TestArmSleepMissionPhrase(t *testing.T) {
	dir := t.TempDir()
	hash, _ := policy.Default().Hash()
	if _, err := EnableGuarded(dir, ArmToken, 8, hash); err != nil {
		t.Fatal(err)
	}
	m := LoadMission(dir)
	if m.Envelope.PolicyHash != hash || m.MissionID == "" || m.SleepState != SleepArmed {
		t.Fatalf("%+v", m)
	}
	if m.Envelope.Digest() == "" {
		t.Fatal("digest")
	}
}

func TestEnvelopeRefusesPolicyChange(t *testing.T) {
	dir := t.TempDir()
	pol := policy.Default()
	hash, _ := pol.Hash()
	if _, err := EnableGuarded(dir, EnableToken, 8, hash); err != nil {
		t.Fatal(err)
	}
	pol.MaxClipUSD = 50
	g := ExecGate{PreviewHash: "0xabc", StartedUnix: timeNowPlus(dir, 1), SessionOK: true, Policy: pol, Now: timeNowPlus(dir, 10)}
	err := AllowHostExecute(dir, g)
	if err == nil || err.Error() != "policy_changed" {
		t.Fatalf("%v", err)
	}
}

func TestEnvelopeRefusesExpired(t *testing.T) {
	pol := policy.Default()
	hash, _ := pol.Hash()
	e := AutonomyEnvelope{PolicyHash: hash, ExpiresAtUnix: 10}
	err := e.RefuseIfStale(pol, ExecGate{Now: 11, SessionOK: true})
	if err == nil || err.Error() != "autonomy_expired" {
		t.Fatalf("%v", err)
	}
}

func TestEnvelopeNamedRefusals(t *testing.T) {
	pol := policy.Default()
	hash, _ := pol.Hash()
	e := AutonomyEnvelope{PolicyHash: hash, MaxOpportunityAgeSec: 5, MaxAutonomyTrades: 1, MaxDailyAutonomyLoss: 10, MaxAutonomyNotional: 10}
	cases := []struct {
		g    ExecGate
		want string
	}{
		{ExecGate{Now: 20, SessionOK: true, Kill: true}, "kill_switch"},
		{ExecGate{Now: 20, SessionOK: false}, "session_expired"},
		{ExecGate{Now: 20, SessionOK: true, MarketUnix: 1}, "stale_market_data"},
		{ExecGate{Now: 20, SessionOK: true, TradesToday: 1}, "max_autonomy_trades"},
		{ExecGate{Now: 20, SessionOK: true, RealizedPnL: -11}, "daily_loss_limit"},
		{ExecGate{Now: 20, SessionOK: true, RequestedUSD: 11}, "max_autonomy_notional"},
		{ExecGate{Now: 20, SessionOK: true, ResearchRequired: true}, "research_unverified"},
		{ExecGate{Now: 20, SessionOK: true, TEERequired: true}, "tee_unverified"},
		{ExecGate{Now: 20, SessionOK: true, ExtraAgentsMissing: true}, "extra_agent_missing"},
		{ExecGate{Now: 20, SessionOK: true, BelowMin: true}, "below_min_notional"},
		{ExecGate{Now: 20, SessionOK: true, InsufficientMargin: true}, "insufficient_margin"},
		{ExecGate{Now: 20, SessionOK: true, WorkspaceID: "b"}, "workspace_mismatch"},
	}
	e.WorkspaceID = "a"
	for _, c := range cases {
		got := e.RefuseIfStale(pol, c.g)
		if got == nil || got.Error() != c.want {
			t.Fatalf("%s: %v", c.want, got)
		}
	}
}

func TestNoTradeEventIsSuccess(t *testing.T) {
	dir := t.TempDir()
	AppendEvent(dir, MissionEvent{Node: "CHALLENGER", Status: "NO-TRADE", Reason: "thesis did not survive challenge", NoTrade: true, Coin: "BTC"})
	log := LoadEvents(dir)
	if log.Challenger != 1 || len(log.Events) != 1 {
		t.Fatalf("%+v", log)
	}
	if !log.Events[0].NoTrade {
		t.Fatal("no-trade")
	}
}

func TestSkillbookHonesty(t *testing.T) {
	dir := t.TempDir()
	b := RecordMemory(dir, KindObservation, "ETH", "host_rank", "sm-1", 1)
	if !strings.HasPrefix(PublicSkillbook(dir)["copy"].(string), "NOT ENOUGH DATA") {
		t.Fatal(PublicSkillbook(dir)["copy"])
	}
	if b.MemoryRoot == "" && LoadSkillbook(dir).MemoryRoot == "" {
		t.Fatal("root")
	}
	raw := LoadSkillbook(dir)
	if strings.Contains(strings.ToLower(raw.Entries[0].Why), "app-sk-") {
		t.Fatal("secret")
	}
}

func TestProofNeverCarriesStrategy(t *testing.T) {
	dir := t.TempDir()
	p := AssembleProof(dir, MissionProof{FillState: "RESTING", OID: "1"})
	if p.FillState != "RESTING" {
		t.Fatal(p.FillState)
	}
	pub := PublicProof(dir)
	if pub["arm"] != false || pub["private_strategy"] != "redacted" {
		t.Fatal(pub)
	}
	rep := StrategyReplay(dir)
	if rep.Live || rep.Trade {
		t.Fatal(rep)
	}
}

func TestSleepStateMap(t *testing.T) {
	if SleepFromStage("scanning", "ACTIVE", "") != SleepWatching {
		t.Fatal("watch")
	}
	if SleepFromStage("researching", "ACTIVE", "") != SleepResearching {
		t.Fatal("research")
	}
	if SleepFromStage("stopped", "STOPPED", "deadline") != SleepStopped {
		t.Fatal("stop")
	}
}

func timeNowPlus(dir string, d int64) int64 {
	m := LoadMission(dir)
	return m.GuardedEnabledUnix + d
}
