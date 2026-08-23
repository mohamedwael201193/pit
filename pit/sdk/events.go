package sdk

import "github.com/mohamedwael201193/pit/internal/phase"

func (c Client) Events() []string {
	return []string{
		phase.Connecting, phase.Authenticating, phase.Sealing, phase.Researching,
		phase.Challenging, phase.AssessingRisk, phase.Scoring, phase.PolicyCheck,
		phase.WaitingForUser, phase.Signing, phase.Submitting, phase.Confirming,
		phase.Executed, phase.Verifying, phase.Resolved, phase.Calibrated, phase.Failed,
	}
}

func (c Client) EventKnown(s string) bool {
	return phase.Known(s)
}
