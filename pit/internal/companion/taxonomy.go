package companion

import "strings"

const (
	TermReadyEligible        = "READY_ELIGIBLE"
	TermReadyStoodDown       = "READY_STOOD_DOWN"
	TermPolicyDenied         = "POLICY_DENIED"
	TermMarketDenied         = "MARKET_DENIED"
	TermCommitteeIncomplete  = "COMMITTEE_INCOMPLETE"
	TermCanceledByUser       = "CANCELED_BY_USER"
	TermDirectNotAuthorized  = "DIRECT_NOT_AUTHORIZED"
	TermDirectCredit         = "DIRECT_CREDIT_INSUFFICIENT"
	TermDirectTimeout        = "DIRECT_PROVIDER_TIMEOUT"
	TermDirectUnavailable     = "DIRECT_PROVIDER_UNAVAILABLE"
	TermDirectRateLimited    = "DIRECT_RATE_LIMITED"
	TermTEESignatureInvalid  = "TEE_SIGNATURE_INVALID"
	TermTEESignerMismatch    = "TEE_SIGNER_MISMATCH"
	TermJobCrashed           = "JOB_CRASHED"
	TermSponsorQuota         = "SPONSOR_QUOTA"
	TermWrongNetwork         = "WRONG_NETWORK"
	TermPollFailed           = "POLL_FAILED"
)

func committeeDenyCode(code string) bool {
	switch strings.ToLower(strings.TrimSpace(code)) {
	case "risk_killed", "challenger_killed", "no_side", "below_min_notional",
		"policy_denied", "kill_switch", "coin_not_allowed", "leverage_above_policy":
		return true
	default:
		return false
	}
}

func policyDenyCode(code string) bool {
	switch strings.ToLower(strings.TrimSpace(code)) {
	case "asset_not_allowed", "kill_switch", "policy_denied", "policy_rejected",
		"leverage", "leverage_above_policy", "daily_loss_halt", "max_trade", "cooldown":
		return true
	default:
		return false
	}
}

func marketDenyCode(code string) bool {
	switch strings.ToLower(strings.TrimSpace(code)) {
	case "empty_envelope", "below_min_notional", "hl_market_unavailable":
		return true
	default:
		return false
	}
}

func namedRolesVerified(roles []map[string]any) bool {
	want := map[string]bool{"researcher": false, "challenger": false, "risk": false}
	for _, rm := range roles {
		role := strings.ToLower(strings.TrimSpace(fmtString(rm["role"])))
		if _, ok := want[role]; !ok {
			continue
		}
		if strings.EqualFold(strings.TrimSpace(fmtString(rm["verify_e2ee"])), "OK") {
			want[role] = true
		}
	}
	return want["researcher"] && want["challenger"] && want["risk"]
}

func okRoleCount(roles []map[string]any) int {
	n := 0
	for _, rm := range roles {
		if strings.EqualFold(strings.TrimSpace(fmtString(rm["verify_e2ee"])), "OK") {
			n++
		}
	}
	return n
}

// TerminalKind is the user-facing terminal taxonomy. POLL_FAILED is never a job terminal.
func TerminalKind(running bool, err, deny string, verify bool, eligible bool, roles []map[string]any) string {
	code := classifyResearch(err)
	if code == TermPollFailed || strings.EqualFold(err, TermPollFailed) {
		if running {
			return ""
		}
	}
	if running {
		return ""
	}
	switch code {
	case "research_cancelled", "CANCELED_BY_USER":
		return TermCanceledByUser
	case TermDirectNotAuthorized, "direct_token_required", "direct_token_expired":
		return TermDirectNotAuthorized
	case TermDirectCredit, "direct_ledger":
		return TermDirectCredit
	case TermDirectTimeout:
		return TermDirectTimeout
	case TermDirectUnavailable:
		return TermDirectUnavailable
	case "direct_rate_limited", "429", TermDirectRateLimited:
		return TermDirectRateLimited
	case "TEE_VERIFY_FAIL", TermTEESignatureInvalid, "TEE_RESPONSE_INVALID":
		return TermTEESignatureInvalid
	case TermTEESignerMismatch, "missing_tee_signer":
		return TermTEESignerMismatch
	case TermJobCrashed:
		return TermJobCrashed
	case TermSponsorQuota:
		return TermSponsorQuota
	case TermWrongNetwork:
		return TermWrongNetwork
	case "POLICY_REJECTED", "asset_not_allowed", "kill_switch":
		return TermPolicyDenied
	case "empty_envelope", "HL_MARKET_UNAVAILABLE":
		return TermMarketDenied
	}
	if policyDenyCode(deny) || policyDenyCode(code) {
		return TermPolicyDenied
	}
	if marketDenyCode(deny) || marketDenyCode(code) {
		return TermMarketDenied
	}
	if committeeDenyCode(deny) || committeeDenyCode(code) {
		if verify || namedRolesVerified(roles) {
			return TermReadyStoodDown
		}
	}
	if verify && eligible {
		return TermReadyEligible
	}
	if verify && !eligible {
		if committeeDenyCode(deny) {
			return TermReadyStoodDown
		}
		if policyDenyCode(deny) {
			return TermPolicyDenied
		}
		if marketDenyCode(deny) {
			return TermMarketDenied
		}
		return TermReadyStoodDown
	}
	if namedRolesVerified(roles) && eligible {
		return TermReadyEligible
	}
	if doneIncomplete(roles, err) {
		return TermCommitteeIncomplete
	}
	if code == "COMMITTEE_INCOMPLETE" {
		return TermCommitteeIncomplete
	}
	if strings.TrimSpace(err) != "" {
		return code
	}
	if len(roles) > 0 && !namedRolesVerified(roles) {
		return TermCommitteeIncomplete
	}
	return ""
}

func doneIncomplete(roles []map[string]any, err string) bool {
	if namedRolesVerified(roles) {
		return false
	}
	if okRoleCount(roles) > 0 && okRoleCount(roles) < 3 {
		return true
	}
	return classifyResearch(err) == "COMMITTEE_INCOMPLETE"
}

func researchCardTitle(kind string) string {
	switch kind {
	case TermReadyEligible:
		return "RESEARCH VERIFIED"
	case TermReadyStoodDown:
		return "COMMITTEE STOOD DOWN"
	case TermCommitteeIncomplete:
		return "RESEARCH INCOMPLETE"
	case TermCanceledByUser:
		return "YOU CANCELLED"
	case TermPolicyDenied:
		return "POLICY BLOCKED THIS"
	case TermMarketDenied:
		return "MARKET NOT USABLE"
	default:
		return "RESEARCH STATUS"
	}
}
