package compute

import (
	"fmt"

	"github.com/mohamedwael201193/pit/internal/config"
	"github.com/mohamedwael201193/pit/internal/deskid"
)

// ProductAsk is the sealed private-book entry. Direct fail stops the operation.
func ProductAsk(net config.Network, deskAuthorized bool, bin string) error {
	if err := deskid.BeforeSealedAsk(deskAuthorized); err != nil {
		return err
	}
	sku := ForNetwork(net)
	if err := SealedAskEnabled(sku); err != nil {
		return err
	}
	if err := RefuseSKUCopy(net, sku.Model); err != nil {
		return err
	}
	if err := DenyRouter(sku.URL); err != nil && sku.URL != "" {
		return err
	}
	if sku.URL == "" {
		return fmt.Errorf("provider_url_required")
	}
	return RunSealedAsk(DirectJob{
		Bin:           bin,
		Role:          Researcher,
		ProviderURL:   sku.URL,
		OnchainSigner: sku.TeeSigner,
		AuthPath:      "auth.json",
		PromptPath:    "prompt.json",
		OutPath:       "out.json",
	})
}
