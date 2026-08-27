package compute

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/mohamedwael201193/pit/internal/config"
	"github.com/mohamedwael201193/pit/internal/deskid"
)

// ProductAsk is the sealed private-book entry. Direct fail stops the operation.
func ProductAsk(net config.Network, deskAuthorized bool, bin string) error {
	return ProductAskEnvelope(net, deskAuthorized, bin, nil, nil)
}

func ProductAskEnvelope(net config.Network, deskAuthorized bool, bin string, publicMarket, privateBook []byte) error {
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
	if err := MustNativeSealer(bin); err != nil {
		return err
	}
	authPath := DirectAuthPath()
	if authPath == "" {
		return fmt.Errorf("direct_token_required")
	}
	raw, err := os.ReadFile(authPath)
	if err != nil {
		return fmt.Errorf("direct_token_required")
	}
	var loaded AuthFile
	if err := json.Unmarshal(raw, &loaded); err != nil {
		return fmt.Errorf("direct_token_required")
	}
	if err := RefuseRouterKey(loaded.Authorization); err != nil {
		return err
	}
	if err := DenyRouter(loaded.URL); err != nil {
		return err
	}
	if !skuURLMatch(sku.URL, loaded.URL) {
		return fmt.Errorf("provider_url_mismatch")
	}
	if len(publicMarket) == 0 || len(privateBook) == 0 {
		return fmt.Errorf("empty_envelope")
	}
	dir, err := os.MkdirTemp("", "pit-ask-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	envelopes, err := Committee(publicMarket, privateBook)
	if err != nil {
		return err
	}
	var jobs []DirectJob
	for _, role := range CommitteeRoles() {
		j, err := MaterializeAsk(dir, sku, role, envelopes[role], loaded.Authorization)
		if err != nil {
			return err
		}
		j.Bin = bin
		jobs = append(jobs, j)
	}
	return RunCommittee(bin, jobs)
}

func skuURLMatch(skuURL, authURL string) bool {
	a := strings.TrimRight(strings.TrimSpace(skuURL), "/")
	b := strings.TrimRight(strings.TrimSpace(authURL), "/")
	return a != "" && b != "" && strings.EqualFold(a, b)
}
