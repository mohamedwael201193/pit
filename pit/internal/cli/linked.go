package cli

import (
	"github.com/mohamedwael201193/pit/internal/config"
	"github.com/mohamedwael201193/pit/internal/hl"
	"github.com/mohamedwael201193/pit/internal/session"
)

func LiveLinked(network, user, workspace, agentAddr string, nowMs int64) (bool, error) {
	net, err := config.ParseNetwork(network)
	if err != nil {
		return false, err
	}
	name, err := session.AgentName(workspace)
	if err != nil {
		return false, err
	}
	raw, err := hl.New(config.For(net)).ExtraAgents(user)
	if err != nil {
		return false, err
	}
	return hl.SessionAgentLinked(raw, name, agentAddr, nowMs), nil
}
