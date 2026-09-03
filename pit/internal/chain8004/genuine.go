package chain8004

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mohamedwael201193/pit/internal/hl"
)

// RecordedFill is the public HYPE fill this desk already proved on Hyperliquid.
const RecordedOID = "531667200134"
const RecordedJob = "4a1d45ec-8c3f-4883-a162-19739accb9cf"
const RecordedResearchTx = "0x1d2113bd683b3ef8be5d74d603018c4bacdd49531bdf201abbc7dea4bb16510b"
const RecordedOrderTx = "0x8c28051bec7bebd7af3b6cc75f7aa034d67f9809f9c30eef9a6c9f84ed6c11fb"

type FillCard struct {
	Type           string `json:"type"`
	AgentRegistry  string `json:"agentRegistry"`
	AgentID        string `json:"agentId"`
	JobID          string `json:"jobId"`
	OID            string `json:"oid"`
	Tag1           string `json:"tag1"`
	Tag2           string `json:"tag2"`
	Value          int    `json:"value"`
	ValueDecimals  int    `json:"valueDecimals"`
	ResearchTx     string `json:"researchTx"`
	OrderTx        string `json:"orderTx"`
}

func ProveFill(raw json.RawMessage, oid string) (hl.FillRow, error) {
	row, ok := hl.FillByOID(raw, oid)
	if !ok {
		return hl.FillRow{}, fmt.Errorf("fill_not_on_venue")
	}
	return row, nil
}

func PublicEvidence(oid string) (job, research, order string) {
	if strings.TrimSpace(oid) == RecordedOID {
		return RecordedJob, RecordedResearchTx, RecordedOrderTx
	}
	return "", "", ""
}

func CanonicalFeedback(agentRegistry string, agentID uint64, oid, tag1, tag2 string) (FillCard, error) {
	if err := RefusePrivateTag(tag1); err != nil {
		return FillCard{}, err
	}
	if err := RefusePrivateTag(tag2); err != nil {
		return FillCard{}, err
	}
	job, research, order := PublicEvidence(oid)
	if job == "" {
		return FillCard{}, fmt.Errorf("unknown_public_job")
	}
	return FillCard{
		Type:          "https://eips.ethereum.org/EIPS/eip-8004#feedback",
		AgentRegistry: agentRegistry,
		AgentID:       fmt.Sprintf("%d", agentID),
		JobID:         job,
		OID:           oid,
		Tag1:          tag1,
		Tag2:          tag2,
		Value:         1,
		ValueDecimals: 0,
		ResearchTx:    research,
		OrderTx:       order,
	}, nil
}

func FeedbackHash(card FillCard) (common.Hash, []byte, error) {
	body, err := json.Marshal(card)
	if err != nil {
		return common.Hash{}, nil, err
	}
	low := strings.ToLower(string(body))
	for _, b := range banned {
		if strings.Contains(low, `"`+b+`"`) {
			return common.Hash{}, nil, fmt.Errorf("private_field_forbidden")
		}
	}
	return common.BytesToHash(crypto.Keccak256(body)), body, nil
}

func DigestSHA256(body []byte) string {
	sum := sha256.Sum256(body)
	return fmt.Sprintf("%x", sum[:])
}

func AgentRegistry(chIdentity string) string {
	return "eip155:16661:" + chIdentity
}
