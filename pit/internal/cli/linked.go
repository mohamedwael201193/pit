package cli

import (
	"github.com/mohamedwael201193/pit/internal/config"
	"github.com/mohamedwael201193/pit/internal/hl"
	"github.com/mohamedwael201193/pit/internal/session"
)

func LiveLinked(network, user, workspace, agentAddr string, nowMs int64) (bool, error) {
	ok, _, err := LiveAgent(network, user, workspace, agentAddr, nowMs)
	return ok, err
}

func LiveAgent(network, user, workspace, agentAddr string, nowMs int64) (bool, int64, error) {
	net, err := config.ParseNetwork(network)
	if err != nil {
		return false, 0, err
	}
	name, err := session.AgentName(workspace)
	if err != nil {
		return false, 0, err
	}
	raw, err := hl.New(config.For(net)).ExtraAgents(user)
	if err != nil {
		return false, 0, err
	}
	ok, until := hl.SessionAgentUntil(raw, name, agentAddr, nowMs)
	return ok, until, nil
}

// LookupAgent is the Hyperliquid PIT-agent list probe. Tests replace it so EnsureLocalSession never mints on a stub.
var LookupAgent = LiveAgent
