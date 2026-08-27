package cli

import (
	"strings"
	"testing"
	"time"

	"github.com/mohamedwael201193/pit/internal/hl"
	"github.com/mohamedwael201193/pit/internal/identity"
	"github.com/mohamedwael201193/pit/internal/policy"
)

func TestSignBoundRecoversSession(t *testing.T) {
	dir := t.TempDir()
	ws := identity.NewWorkspaceID()
	sf, err := CreateLocalSession(dir, ws, "testnet", policy.Default().Version)
	if err != nil {
		t.Fatal(err)
	}
	live, err := LiveFromDisk(dir, false, time.Now().UnixMilli())
	if err != nil {
		t.Fatal(err)
	}
	raw, err := hl.BuildOrder(1, true, "2500", "0.004", "0x11111111111111111111111111111111")
	if err != nil {
		t.Fatal(err)
	}
	env, err := SignBound(dir, live, "testnet", raw, 1700000000000)
	if err != nil {
		t.Fatal(err)
	}
	got, err := hl.RecoverL1(env, false)
	if err != nil || !strings.EqualFold(got, sf.AgentAddr) {
		t.Fatalf("%s %s %v", got, sf.AgentAddr, err)
	}
}
