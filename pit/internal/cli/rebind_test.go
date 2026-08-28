package cli

import (
	"encoding/json"
	"testing"

	"github.com/mohamedwael201193/pit/internal/engine"
)

func TestReportFromRoleMapsChallengerKilled(t *testing.T) {
	roles := []map[string]any{
		{"role": "researcher", "verify_e2ee": "OK", "proposed_side": "buy"},
		{"role": "challenger", "verify_e2ee": "OK", "survives": false, "kill": false},
		{"role": "risk", "verify_e2ee": "OK", "survives": true, "kill": false},
	}
	rep, err := reportFromRoleMaps(roles)
	if err != nil {
		t.Fatal(err)
	}
	got := engine.EvaluateCommittee(2500, 4, 15, 15, "", "ETH", []string{"ETH"}, 1, 1, false, rep.Researcher, rep.Challenger, rep.Risk)
	if got.Eligible || got.Deny != "challenger_killed" {
		t.Fatalf("%+v", got)
	}
	var ch map[string]any
	if json.Unmarshal(rep.Challenger, &ch) != nil || ch["survives"] != false {
		t.Fatalf("%s", rep.Challenger)
	}
}
