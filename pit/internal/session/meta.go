package session

import (
	"fmt"

	"github.com/google/uuid"
)

type Meta struct {
	ID        string
	AgentAddr string
	Workspace string
	Network   string
	PolicyVer string
	Expires   int64
}

func NewID() string {
	return uuid.NewString()
}

func (m Meta) Session() Session {
	return Session{
		ID:        m.ID,
		Workspace: m.Workspace,
		AgentAddr: m.AgentAddr,
		Expires:   m.Expires,
		PolicyVer: m.PolicyVer,
		Network:   m.Network,
	}
}

func (m Meta) Validate() error {
	if m.ID == "" || m.AgentAddr == "" || m.Workspace == "" {
		return fmt.Errorf("incomplete_session")
	}
	return nil
}
