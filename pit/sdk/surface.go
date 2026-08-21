package sdk

import (
	"github.com/mohamedwael201193/pit/internal/config"
	"github.com/mohamedwael201193/pit/internal/phase"
	"github.com/mohamedwael201193/pit/internal/policy"
	"github.com/mohamedwael201193/pit/internal/watch"
)

func (c Client) Phases() []string {
	return []string{
		phase.Connecting, phase.Sealing, phase.Researching, phase.Challenging,
		phase.AssessingRisk, phase.PolicyCheck, phase.WaitingForUser, phase.Executed,
		phase.Verifying, phase.Calibrated, phase.Failed,
	}
}

func (c Client) EmptyWatch() string {
	return watch.Attention(0)
}

func (c Client) NetworkChain() config.Chain {
	return config.For(c.Network)
}

func (c Client) PolicyCards() policy.Policy {
	return policy.Default()
}
