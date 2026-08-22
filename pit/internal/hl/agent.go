package hl

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/mohamedwael201193/pit/internal/session"
)

// ApproveAgent is signed by the master wallet, never by the session key.
type ApproveAgent struct {
	Type         string `json:"type"`
	AgentAddress string `json:"agentAddress"`
	AgentName    string `json:"agentName"`
	Nonce        int64  `json:"nonce"`
}

func BuildApproveAgent(agentAddr, workspaceID string, now time.Time) (json.RawMessage, error) {
	if agentAddr == "" {
		return nil, fmt.Errorf("agent_required")
	}
	name, err := session.AgentName(workspaceID)
	if err != nil {
		return nil, err
	}
	if strings.Contains(strings.ToLower(agentAddr), "session") {
		return nil, fmt.Errorf("session_cannot_approve")
	}
	a := ApproveAgent{
		Type:         "approveAgent",
		AgentAddress: agentAddr,
		AgentName:    name,
		Nonce:        now.UnixMilli(),
	}
	return json.Marshal(a)
}

func SessionMustNotSign(actionType string) error {
	if actionType == "approveAgent" {
		return fmt.Errorf("master_wallet_only")
	}
	return nil
}
