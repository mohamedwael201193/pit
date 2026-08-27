package compute

import (
	"encoding/json"
	"fmt"
)

type PrivateBook struct {
	Wallet     string   `json:"wallet"`
	Workspace  string   `json:"workspace"`
	Network    string   `json:"network"`
	PolicyHash string   `json:"policyHash"`
	Positions  []string `json:"positions"`
	Note       string   `json:"note"`
}

func BuildPrivateBook(wallet, workspace, network, policyHash string) ([]byte, error) {
	if wallet == "" || workspace == "" || network == "" {
		return nil, fmt.Errorf("empty_envelope")
	}
	b := PrivateBook{
		Wallet:     wallet,
		Workspace:  workspace,
		Network:    network,
		PolicyHash: policyHash,
		Positions:  []string{},
		Note:       "No local fills are recorded on this machine yet. PIT does not invent positions.",
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
