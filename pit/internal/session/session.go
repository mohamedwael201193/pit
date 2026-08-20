package session

import (
	"fmt"
	"time"
)

var AllowedActions = map[string]struct{}{
	"order":  {},
	"cancel": {},
}

var DeniedActions = map[string]struct{}{
	"withdraw3": {}, "usdSend": {}, "spotSend": {}, "usdClassTransfer": {},
	"sendAsset": {}, "updateLeverage": {}, "updateIsolatedMargin": {},
	"vaultTransfer": {}, "subAccountTransfer": {}, "createSubAccount": {},
	"subAccountModify": {}, "approveAgent": {}, "approveBuilderFee": {},
	"cDeposit": {}, "cWithdraw": {}, "tokenDelegate": {}, "spotDeploy": {},
	"convertToMultiSigUser": {}, "userSetAbstraction": {}, "userDexAbstraction": {},
	"twapOrder": {}, "twapCancel": {}, "batchModify": {}, "modify": {},
	"scheduleCancel": {}, "setReferrer": {}, "reserveRequestWeight": {}, "noop": {},
}

type Denial struct {
	Action string
	Reason string
}

func (d Denial) Error() string {
	return fmt.Sprintf("PIT_DENY action=%s reason=%s", d.Action, d.Reason)
}

func CheckAction(actionType string) error {
	if actionType == "" {
		return Denial{Action: actionType, Reason: "empty_action"}
	}
	if _, ok := AllowedActions[actionType]; ok {
		return nil
	}
	reason := "not_allowlisted"
	if _, listed := DeniedActions[actionType]; listed {
		reason = "explicit_deny"
	}
	return Denial{Action: actionType, Reason: reason}
}

type Session struct {
	ID        string
	Workspace string
	AgentAddr string
	Expires   int64
	Revoked   bool
	Kill      bool
	PolicyVer string
	Network   string
}

func CheckSession(s Session, nowMs int64, wantPolicy, wantWorkspace string) error {
	if s.ID == "" || s.AgentAddr == "" {
		return fmt.Errorf("PIT_DENY reason=unknown_session")
	}
	if wantWorkspace != "" && s.Workspace != wantWorkspace {
		return fmt.Errorf("PIT_DENY reason=wrong_workspace")
	}
	if s.Revoked {
		return fmt.Errorf("PIT_DENY reason=revoked_session")
	}
	if s.Kill {
		return fmt.Errorf("PIT_DENY reason=kill_switch")
	}
	if nowMs >= s.Expires {
		return fmt.Errorf("PIT_DENY reason=expired_session")
	}
	if wantPolicy != "" && s.PolicyVer != wantPolicy {
		return fmt.Errorf("PIT_DENY reason=policy_version_mismatch")
	}
	return nil
}

func TTL(now time.Time, hours int) int64 {
	if hours <= 0 {
		hours = 1
	}
	return now.Add(time.Duration(hours) * time.Hour).UnixMilli()
}
