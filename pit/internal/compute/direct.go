package compute

import (
	"fmt"
	"strings"
)

type DirectJob struct {
	Bin           string
	AuthPath      string
	PromptPath    string
	OutPath       string
	Role          Role
	ProviderURL   string
	OnchainSigner string
}

func PrepareDirect(j DirectJob) ([]string, error) {
	if j.Bin == "" {
		return nil, fmt.Errorf("sealer_not_wired")
	}
	if strings.HasSuffix(strings.ToLower(j.Bin), ".py") {
		return nil, fmt.Errorf("python_gate_not_product")
	}
	if err := DenyRouter(j.ProviderURL); err != nil {
		return nil, err
	}
	if j.AuthPath == "" || j.PromptPath == "" || j.OutPath == "" {
		return nil, fmt.Errorf("incomplete_direct_job")
	}
	switch j.Role {
	case Researcher, Challenger, Risk:
	default:
		return nil, fmt.Errorf("bad_role")
	}
	if j.OnchainSigner == "" {
		return nil, fmt.Errorf("missing_tee_signer")
	}
	return []string{
		"-auth", j.AuthPath,
		"-prompt", j.PromptPath,
		"-out", j.OutPath,
		"-role", string(j.Role),
	}, nil
}

