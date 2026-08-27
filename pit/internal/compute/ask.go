package compute

import (
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
	return ProductAskAuth(net, deskAuthorized, bin, publicMarket, privateBook, AuthFile{})
}

func ProductAskAuth(net config.Network, deskAuthorized bool, bin string, publicMarket, privateBook []byte, loaded AuthFile) error {
	_, err := ProductAskReport(net, deskAuthorized, bin, publicMarket, privateBook, loaded)
	return err
}

type AskReport struct {
	Roles []map[string]any `json:"roles"`
	Note  string           `json:"note"`
}

func ProductAskReport(net config.Network, deskAuthorized bool, bin string, publicMarket, privateBook []byte, loaded AuthFile) (AskReport, error) {
	if err := deskid.BeforeSealedAsk(deskAuthorized); err != nil {
		return AskReport{}, err
	}
	sku := ForNetwork(net)
	if err := SealedAskEnabled(sku); err != nil {
		return AskReport{}, err
	}
	if err := RefuseSKUCopy(net, sku.Model); err != nil {
		return AskReport{}, err
	}
	if err := DenyRouter(sku.URL); err != nil && sku.URL != "" {
		return AskReport{}, err
	}
	if sku.URL == "" {
		return AskReport{}, fmt.Errorf("provider_url_required")
	}
	if err := MustNativeSealer(bin); err != nil {
		return AskReport{}, err
	}
	if strings.TrimSpace(loaded.Authorization) == "" {
		var err error
		loaded, err = LoadEnvAuthFile()
		if err != nil {
			return AskReport{}, err
		}
	}
	if err := RefuseRouterKey(loaded.Authorization); err != nil {
		return AskReport{}, err
	}
	if err := DenyRouter(loaded.URL); err != nil {
		return AskReport{}, err
	}
	if !skuURLMatch(sku.URL, loaded.URL) {
		return AskReport{}, fmt.Errorf("provider_url_mismatch")
	}
	if len(publicMarket) == 0 || len(privateBook) == 0 {
		return AskReport{}, fmt.Errorf("empty_envelope")
	}
	dir, err := os.MkdirTemp("", "pit-ask-")
	if err != nil {
		return AskReport{}, err
	}
	defer os.RemoveAll(dir)
	envelopes, err := Committee(publicMarket, privateBook)
	if err != nil {
		return AskReport{}, err
	}
	var jobs []DirectJob
	for _, role := range CommitteeRoles() {
		j, err := MaterializeAsk(dir, sku, role, envelopes[role], loaded.Authorization)
		if err != nil {
			return AskReport{}, err
		}
		j.Bin = bin
		jobs = append(jobs, j)
	}
	if err := RunCommittee(bin, jobs); err != nil {
		return AskReport{}, err
	}
	rep := AskReport{Note: HonestLabel(IndependenceNote()), Roles: make([]map[string]any, 0, len(jobs))}
	for _, j := range jobs {
		rep.Roles = append(rep.Roles, PublicRoleEvidence(j))
	}
	return rep, nil
}

func skuURLMatch(skuURL, authURL string) bool {
	a := strings.TrimRight(strings.TrimSpace(skuURL), "/")
	b := strings.TrimRight(strings.TrimSpace(authURL), "/")
	return a != "" && b != "" && strings.EqualFold(a, b)
}
