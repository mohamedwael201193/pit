package sdk

import (
	"github.com/mohamedwael201193/pit/internal/config"
	"github.com/mohamedwael201193/pit/internal/policy"
)

// Client is the typed surface for apps that must not hold a session secret.
type Client struct {
	Network config.Network
}

type Status struct {
	Network   config.Network `json:"network"`
	Workspace string         `json:"workspace,omitempty"`
	CanSign   bool           `json:"canSign"`
}

func (c Client) Status() Status {
	return Status{Network: c.Network, CanSign: false}
}

func (c Client) DefaultPolicy() policy.Policy {
	return policy.Default()
}

func (c Client) CanHoldSession() bool {
	return false
}

func (c Client) Explorer(tx string) string {
	ch := config.For(c.Network)
	if tx == "" {
		return ch.Explorer
	}
	return ch.Explorer + "/tx/" + tx
}
