package chain8004

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mohamedwael201193/pit/internal/chainrpc"
	"github.com/mohamedwael201193/pit/internal/config"
	"github.com/mohamedwael201193/pit/internal/identity"
)

func ReporterKey() string {
	if v := strings.TrimSpace(os.Getenv("PIT_8004_REPORTER_KEY")); v != "" {
		return v
	}
	return strings.TrimSpace(os.Getenv("PIT_DEPLOYER_KEY"))
}

func refuseSessionSigner(addr common.Address) error {
	if strings.EqualFold(addr.Hex(), RecordedHLAgent) {
		return fmt.Errorf("session_agent_cannot_report")
	}
	return nil
}

// SubmitFillFeedback posts 8-arg giveFeedback from a reporter EOA that is not
// the identity owner and not the Hyperliquid session agent. The fill must be
// present on Hyperliquid user_fills. Idempotent if this reporter already wrote.
func SubmitFillFeedback(ch config.Chain, reporterKey string, agentID uint64, oid string, fills json.RawMessage) (string, error) {
	if _, err := ProveFill(fills, oid); err != nil {
		return "", err
	}
	reporter, err := chainrpc.AddressOfKey(reporterKey)
	if err != nil {
		return "", err
	}
	if err := refuseSessionSigner(reporter); err != nil {
		return "", err
	}
	owner, err := OwnerOf(ch, agentID)
	if err != nil {
		return "", err
	}
	own, err := identity.NormalizeAddress(owner)
	if err != nil {
		return "", err
	}
	rep, err := identity.NormalizeAddress(reporter.Hex())
	if err != nil {
		return "", err
	}
	if err := FeedbackAllowed(own, rep, rep); err != nil {
		return "", err
	}
	idx, err := LastIndex(ch, agentID, reporter.Hex())
	if err != nil {
		return "", err
	}
	if idx > 0 {
		fb, rerr := ReadFeedback(ch, agentID, reporter.Hex(), idx)
		if rerr == nil && fb.Tag1 == "hype_fill" && fb.Tag2 == "successful_job" {
			return "", fmt.Errorf("already_reported")
		}
	}
	card, err := CanonicalFeedback(AgentRegistry(ch.Identity8004), agentID, oid, "hype_fill", "successful_job")
	if err != nil {
		return "", err
	}
	hash, _, err := FeedbackHash(card)
	if err != nil {
		return "", err
	}
	data, err := EncodeGiveFeedback(FeedbackArgs{
		AgentID: agentID, Value: 1, Tag1: card.Tag1, Tag2: card.Tag2,
		Endpoint: WebEndpoint, FeedbackURI: ProofURL, FeedbackHash: hash,
	})
	if err != nil {
		return "", err
	}
	return chainrpc.Send(ch.RPC, ch.ChainID, reporterKey, ch.Reputation8004, data)
}

func SetURIFromOwner(ch config.Chain, ownerKey string, agentID uint64, uri string) (string, error) {
	from, err := chainrpc.AddressOfKey(ownerKey)
	if err != nil {
		return "", err
	}
	if err := refuseSessionSigner(from); err != nil {
		return "", err
	}
	owner, err := OwnerOf(ch, agentID)
	if err != nil {
		return "", err
	}
	if !strings.EqualFold(owner, from.Hex()) {
		return "", fmt.Errorf("not_identity_owner")
	}
	cur, _ := AgentURI(ch, agentID)
	if strings.TrimSpace(cur) == strings.TrimSpace(uri) {
		return "", fmt.Errorf("uri_unchanged")
	}
	data, err := EncodeSetAgentURI(agentID, uri)
	if err != nil {
		return "", err
	}
	return chainrpc.Send(ch.RPC, ch.ChainID, ownerKey, ch.Identity8004, data)
}
