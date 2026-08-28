package compute

import (
	"encoding/json"
	"fmt"
	"strings"
)

type PrivateBook struct {
	Wallet     string   `json:"wallet"`
	Workspace  string   `json:"workspace"`
	Network    string   `json:"network"`
	PolicyHash string   `json:"policyHash"`
	Positions  []string `json:"positions"`
	Hypothesis string   `json:"hypothesis"`
	Note       string   `json:"note"`
}

func ParseHypothesis(s string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "", "none", "flat", "stand", "stand-down", "stand_down", "no-bias", "nobias":
		return "none", nil
	case "long", "buy":
		return "long", nil
	case "short", "sell":
		return "short", nil
	default:
		return "", fmt.Errorf("bad_hypothesis")
	}
}

func BuildPrivateBook(wallet, workspace, network, policyHash string) ([]byte, error) {
	return BuildPrivateBookHypothesis(wallet, workspace, network, policyHash, "none")
}

func BuildPrivateBookHypothesis(wallet, workspace, network, policyHash, hypothesis string) ([]byte, error) {
	if wallet == "" || workspace == "" || network == "" {
		return nil, fmt.Errorf("empty_envelope")
	}
	hyp, err := ParseHypothesis(hypothesis)
	if err != nil {
		return nil, err
	}
	note := "No local fills are recorded on this machine yet. PIT does not invent positions."
	switch hyp {
	case "long":
		note = "User sealed hypothesis: considering a long. Committee may reject with none. Host still sizes. PIT does not invent fills."
	case "short":
		note = "User sealed hypothesis: considering a short. Committee may reject with none. Host still sizes. PIT does not invent fills."
	}
	b := PrivateBook{
		Wallet:     wallet,
		Workspace:  workspace,
		Network:    network,
		PolicyHash: policyHash,
		Positions:  []string{},
		Hypothesis: hyp,
		Note:       note,
	}
	return json.Marshal(b)
}

func MarketJSON(v any) ([]byte, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	if len(b) == 0 || string(b) == "null" {
		return nil, fmt.Errorf("empty_envelope")
	}
	return b, nil
}
