package compute

import (
	"fmt"
	"strings"

	"github.com/mohamedwael201193/pit/internal/config"
)

func MatchProvider(n config.Network, provider, model, signer string) error {
	sku := ForNetwork(n)
	if !strings.EqualFold(strings.TrimSpace(provider), sku.Provider) {
		return fmt.Errorf("provider_spoof")
	}
	if strings.TrimSpace(model) != sku.Model {
		return fmt.Errorf("model_spoof")
	}
	if !strings.EqualFold(strings.TrimSpace(signer), sku.TeeSigner) {
		return fmt.Errorf("signer_spoof")
	}
	return nil
}
