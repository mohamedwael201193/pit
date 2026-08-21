package phase

// Named product states. UI may only advance when the host emits these.
const (
	Connecting      = "CONNECTING"
	Authenticating  = "AUTHENTICATING"
	Sealing         = "SEALING"
	Researching     = "RESEARCHING"
	Challenging     = "CHALLENGING"
	AssessingRisk   = "ASSESSING_RISK"
	Scoring         = "SCORING"
	PolicyCheck     = "POLICY_CHECK"
	WaitingForUser  = "WAITING_FOR_USER"
	Signing         = "SIGNING"
	Submitting      = "SUBMITTING"
	Confirming      = "CONFIRMING"
	Executed        = "EXECUTED"
	Verifying       = "VERIFYING"
	Resolved        = "RESOLVED"
	Calibrated      = "CALIBRATED"
	Failed          = "FAILED"
)

func Known(s string) bool {
	switch s {
	case Connecting, Authenticating, Sealing, Researching, Challenging, AssessingRisk,
		Scoring, PolicyCheck, WaitingForUser, Signing, Submitting, Confirming,
		Executed, Verifying, Resolved, Calibrated, Failed:
		return true
	}
	return false
}
