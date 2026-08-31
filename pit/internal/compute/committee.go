package compute

import (
	"encoding/json"
	"fmt"
)

// Progress maps to real backend steps. UI must not animate these on a timer.
type Progress string

const (
	ProgBook    Progress = "PRIVATE_BOOK"
	ProgSeal    Progress = "SEALING"
	ProgTEE     Progress = "TEE"
	ProgSig     Progress = "TEE_SIGNATURE"
	ProgSigner  Progress = "ONCHAIN_SIGNER"
	ProgStorage Progress = "STORAGE"
	ProgReceipt Progress = "RECEIPT"
	ProgCalib   Progress = "CALIBRATION"
)

type Event struct {
	Progress Progress
	OK       bool
	Detail   string
}

type Role string

const (
	Researcher Role = "researcher"
	Challenger Role = "challenger"
	Risk       Role = "risk"
)

func Next(cur Progress) (Progress, error) {
	order := []Progress{ProgBook, ProgSeal, ProgTEE, ProgSig, ProgSigner, ProgStorage, ProgReceipt, ProgCalib}
	for i, p := range order {
		if p == cur && i+1 < len(order) {
			return order[i+1], nil
		}
	}
	return "", fmt.Errorf("no_next")
}

func RoleInstruction(role Role) string {
	switch role {
	case Researcher:
		return `ROLE=researcher. Read live market facts and the sealed hypothesis (none|long|short). Reply with JSON only: {"proposed_side":"buy"|"sell"|"none"}. If hypothesis is long, confirm buy unless live facts clearly contradict a long. If hypothesis is short, confirm sell unless live facts clearly contradict a short. If hypothesis is none, pick buy, sell, or none from the live market facts; prefer a side when funding, mark/oracle gap, open interest, or momentum in this envelope leans one way. Do not echo none just because hypothesis is none. Never output size, leverage, withdraw, transfer, or permissions. The host sizes.`
	case Challenger:
		return `ROLE=challenger. Reply with JSON only: {"survives":true|false,"kill":false}. Challenge researcher_thesis in this envelope against the live market facts. If proposed_side is none or missing, survives=false. If proposed_side is buy or sell, survives=true only when that exact side is still justified by the live facts; otherwise survives=false. Never size.`
	case Risk:
		return `ROLE=risk. Reply with JSON only: {"kill":false,"survives":true}. Read researcher_thesis. Kill if that side is unsafe on this book. Never size.`
	default:
		return ""
	}
}

func Envelope(role Role, publicMarket, privateBook []byte) ([]byte, error) {
	switch role {
	case Researcher, Challenger, Risk:
	default:
		return nil, fmt.Errorf("bad_role")
	}
	if len(publicMarket) == 0 || len(privateBook) == 0 {
		return nil, fmt.Errorf("empty_envelope")
	}
	inst := RoleInstruction(role)
	out := make([]byte, 0, len(inst)+16+len(publicMarket)+len(privateBook))
	out = append(out, []byte(inst+"|"+string(role)+"|")...)
	out = append(out, publicMarket...)
	out = append(out, '|')
	out = append(out, privateBook...)
	return out, nil
}

// WithResearcherThesis copies the researcher's JSON into the public market so
// challenger and risk challenge that exact side instead of an empty hypothesis.
func WithResearcherThesis(publicMarket, prior []byte) []byte {
	if len(publicMarket) == 0 {
		return publicMarket
	}
	if len(prior) == 0 {
		return publicMarket
	}
	var m map[string]any
	if err := json.Unmarshal(publicMarket, &m); err != nil {
		return publicMarket
	}
	var thesis any
	if err := json.Unmarshal(prior, &thesis); err != nil {
		m["researcher_thesis"] = string(prior)
	} else {
		m["researcher_thesis"] = thesis
	}
	b, err := json.Marshal(m)
	if err != nil {
		return publicMarket
	}
	return b
}
