package companion

import "testing"

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

func TestStandDownIsNotCrash(t *testing.T) {
	roles := []map[string]any{
		{"role": "researcher", "verify_e2ee": "OK"},
		{"role": "challenger", "verify_e2ee": "OK"},
		{"role": "risk", "verify_e2ee": "OK"},
	}
	if TerminalKind(false, "risk_killed", "risk_killed", true, false, roles) != TermReadyStoodDown {
		t.Fatal(TerminalKind(false, "risk_killed", "risk_killed", true, false, roles))
	}
	if researchCardTitle(TermReadyStoodDown) == "RESEARCH VERIFIED" {
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
	if researchCardTitle(TermReadyEligible) != "RESEARCH VERIFIED" {
		t.Fatal("title")
	}
	if researchCardTitle(TermCommitteeIncomplete) == "RESEARCH VERIFIED" {
		t.Fatal("incomplete")
	}
}
