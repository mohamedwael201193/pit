package auto

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mohamedwael201193/pit/internal/policy"
)

func TestEnableGuardedRequiresPhrase(t *testing.T) {
	dir := t.TempDir()
	if _, err := EnableGuarded(dir, "please", 24, "abc"); err == nil {
		t.Fatal("phrase")
	}
	if _, err := SetMode(dir, ModeGuarded); err == nil {
		t.Fatal("set mode")
	}
	hash, _ := policy.Default().Hash()
	m, err := EnableGuarded(dir, EnableToken, 24, hash)
	if err != nil || m.Mode != ModeGuarded || !m.Running {
		t.Fatalf("%+v %v", m, err)
	}
	p := Load(dir)
	if p.Execute || !p.AutoResearch {
		t.Fatalf("%+v", p)
	}
}

func TestGuardedRunningRecoversFromDisk(t *testing.T) {
	dir := t.TempDir()
	hash, _ := policy.Default().Hash()
	if _, err := EnableGuarded(dir, EnableToken, 8, hash); err != nil {
		t.Fatal(err)
	}
	got := LoadMission(dir)
	if !got.Running || got.Mode != ModeGuarded {
		t.Fatalf("%+v", got)
	}
	again := LoadMission(dir)
	if !again.Running || again.GuardedUntilUnix != got.GuardedUntilUnix {
		t.Fatal("lost mission")
	}
}

func TestStopRecovers(t *testing.T) {
	dir := t.TempDir()
	hash, _ := policy.Default().Hash()
	if _, err := EnableGuarded(dir, EnableToken, 1, hash); err != nil {
		t.Fatal(err)
	}
	Stop(dir, "kill_switch")
	m := LoadMission(dir)
	if m.Mode != ModeManual || m.Running || m.LastStop != "kill_switch" {
		t.Fatalf("%+v", m)
	}
	raw, _ := os.ReadFile(filepath.Join(dir, "mission.json"))
	if string(raw) == "" {
		t.Fatal("persist")
	}
	got := LoadMission(dir)
	if got.LastStop != "kill_switch" {
		t.Fatal("recover")
	}
}

func TestAllowHostExecuteGates(t *testing.T) {
	dir := t.TempDir()
	pol := policy.Default()
	hash, _ := pol.Hash()
	g := ExecGate{PreviewHash: "0x1", StartedUnix: 1, SessionOK: true, Policy: pol, Now: 2}
	if err := AllowHostExecute(dir, g); err == nil {
		t.Fatal("manual must refuse")
	}
	if _, err := EnableGuarded(dir, EnableToken, 24, hash); err != nil {
		t.Fatal(err)
	}
	m := LoadMission(dir)
	g.StartedUnix = m.GuardedEnabledUnix + 1
	g.Now = m.GuardedEnabledUnix + 10
	g.OpenCount = 1
	if err := AllowHostExecute(dir, g); err == nil {
		t.Fatal("open positions")
	}
	g.OpenCount = 0
	if err := AllowHostExecute(dir, g); err != nil {
		t.Fatal(err)
	}
	RecordAction(dir, "posted", "ETH", "0x1", "99", "")
	if err := AllowHostExecute(dir, g); err == nil {
		t.Fatal("duplicate")
	}
}

func TestResearchOnlyNeverExecutes(t *testing.T) {
	dir := t.TempDir()
	if _, err := SetMode(dir, ModeResearch); err != nil {
		t.Fatal(err)
	}
	g := ExecGate{PreviewHash: "0x2", SessionOK: true, Policy: policy.Default(), Now: 10, StartedUnix: 9}
	if err := AllowHostExecute(dir, g); err == nil {
		t.Fatal("research only")
	}
}

func TestStopReasonDeadline(t *testing.T) {
	m := Mission{Mode: ModeGuarded, GuardedUntilUnix: 10, DeadlineUnix: 10}
	if StopReason(m, 11, false, true, 0, 0, policy.Default()) != "autonomy_expired" {
		t.Fatal("deadline")
	}
	if StopReason(m, 5, true, true, 0, 0, policy.Default()) != "kill_switch" {
		t.Fatal("kill")
	}
}

func TestHaltVsExecBlock(t *testing.T) {
	m := Mission{Mode: ModeGuarded, Running: true, GuardedUntilUnix: 4_000_000_000}
	pol := policy.Default()
	if MissionHaltReason(m, 10, false, true, 0, pol) != "" {
		t.Fatal("halt")
	}
	if ExecBlockReason(1, pol) != "max_open_positions" {
		t.Fatal("block")
	}
	if StopReason(m, 10, false, true, 1, 0, pol) != "" {
		t.Fatal("open positions must not halt the mission")
	}
	dir := t.TempDir()
	hash, _ := pol.Hash()
	if _, err := EnableGuarded(dir, EnableToken, 8, hash); err != nil {
		t.Fatal(err)
	}
	got := LoadMission(dir)
	if got.Stage != "starting" || got.NextScanUnix == 0 {
		t.Fatalf("%+v", got)
	}
	p := Load(dir)
	if p.LastScanUnix != 0 || p.LastResearchCoin != "" {
		t.Fatalf("%+v", p)
	}
}

func TestLifeAndPublicStatus(t *testing.T) {
	if Life(Mission{Mode: ModeManual}, false, 1) != "READY" {
		t.Fatal("ready")
	}
	if Life(Mission{Mode: ModeGuarded, Running: true, GuardedUntilUnix: 100}, false, 1) != "ACTIVE" {
		t.Fatal("active")
	}
	if Life(Mission{Mode: ModeGuarded, Running: true, GuardedUntilUnix: 10}, true, 1) != "BLOCKED" {
		t.Fatal("blocked")
	}
	dir := t.TempDir()
	p := Public(dir)
	if p["status"] != "READY" || p["execute"] == true {
		t.Fatalf("%+v", p)
	}
}

func TestRecordBlockDoesNotHalt(t *testing.T) {
	dir := t.TempDir()
	hash, _ := policy.Default().Hash()
	if _, err := EnableGuarded(dir, EnableToken, 8, hash); err != nil {
		t.Fatal(err)
	}
	RecordBlock(dir, "max_open_positions", "BTC")
	m := LoadMission(dir)
	if !m.Running || m.Mode != ModeGuarded || m.LastStop != "" {
		t.Fatalf("%+v", m)
	}
	if Life(m, false, m.GuardedUntilUnix-1) != "ACTIVE" {
		t.Fatal(Life(m, false, m.GuardedUntilUnix-1))
	}
	if Explain("max_open_positions") == "" {
		t.Fatal("explain")
	}
}
