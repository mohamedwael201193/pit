package compute

import (
	"fmt"
	"strings"

	"github.com/mohamedwael201193/pit/internal/config"
)

// SealJob is the private-book path. Router URLs are rejected. Missing sealer stops the op.
type SealJob struct {
	Network       config.Network
	Role          Role
	PublicMarket  []byte
	PrivateBook   []byte
	Scheme        string
	ProviderURL   string
	TeeSigner     string
	SealerBin     string
}

func BindSeal(j SealJob) (DirectJob, []byte, error) {
	sku := ForNetwork(j.Network)
	if err := SealedAskEnabled(sku); err != nil {
		return DirectJob{}, nil, err
	}
	if err := RefuseSKUCopy(j.Network, sku.Model); err != nil {
		return DirectJob{}, nil, err
	}
	if strings.TrimSpace(j.Scheme) != SchemeE2EE {
		return DirectJob{}, nil, fmt.Errorf("TEE_VERIFY_FAIL")
	}
	if !strings.EqualFold(strings.TrimSpace(j.TeeSigner), sku.TeeSigner) {
		return DirectJob{}, nil, fmt.Errorf("tee_signer_mismatch")
	}
	url := j.ProviderURL
	if url == "" {
		url = sku.URL
	}
	if err := DenyRouter(url); err != nil {
		return DirectJob{}, nil, err
	}
	env, err := Envelope(j.Role, j.PublicMarket, j.PrivateBook)
	if err != nil {
		return DirectJob{}, nil, err
	}
	bin := j.SealerBin
	if bin == "" {
		bin = LookBin()
	}
	dj := DirectJob{
		Bin:           bin,
		Role:          j.Role,
		ProviderURL:   url,
		OnchainSigner: sku.TeeSigner,
	}
	if err := MustNativeSealer(dj.Bin); err != nil {
		return DirectJob{}, nil, err
	}
	return dj, env, nil
}
