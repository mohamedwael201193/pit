package compute

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/mohamedwael201193/pit/internal/config"
)

const OfficialCatalogURL = "https://router-api.0g.ai/v1/models"

type CatalogEntry struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Type          string `json:"type"`
	Verifiability string `json:"verifiability"`
	TeeAttested   bool   `json:"tee_attested,omitempty"`
	ProviderCount int    `json:"provider_count,omitempty"`
	PrivateBook   bool   `json:"private_book"`
	Path          string `json:"path"`
	UsableFor     string `json:"usable_for"`
	Note          string `json:"note"`
}

func classifyCatalog(id, ver string, provenDirect string) CatalogEntry {
	v := strings.TrimSpace(ver)
	if v == "" {
		v = "unproven"
	}
	e := CatalogEntry{
		ID:            id,
		Verifiability: v,
		Path:          "catalog-listing",
		PrivateBook:   false,
		UsableFor:     "not used for inference by PIT",
		Note:          "Official catalog listing only. PIT never sends the private book or desk commands to Router.",
	}
	if strings.EqualFold(v, "TeeML") && strings.EqualFold(id, provenDirect) {
		e.Note = "This name matches the Direct TeeML SKU this workspace already uses for sealed research. Catalog presence is not a second path."
	}
	if strings.EqualFold(v, "TeeTLS") {
		e.UsableFor = "listed TeeTLS — not Direct TeeML on this workspace"
		e.Note = "TeeTLS is a verifiable relay, not private-book TeeML. PIT will not call it and will not fall back to Router."
	}
	if strings.EqualFold(v, "unproven") || v == "" {
		e.Verifiability = "unproven"
		e.UsableFor = "not used"
		e.Note = "No TeeML/TeeTLS flag on the official listing. PIT will not call it."
	}
	return e
}

func ListOfficialCatalog(ctx context.Context) ([]CatalogEntry, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 8*time.Second)
		defer cancel()
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, OfficialCatalogURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "pit/1.0")
	client := &http.Client{Timeout: 8 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("catalog_http_%d", resp.StatusCode)
	}
	var pack struct {
		Data []struct {
			ID            string `json:"id"`
			Name          string `json:"name"`
			Type          string `json:"type"`
			Verifiability string `json:"verifiability"`
			TeeAttested   bool   `json:"tee_attested"`
			ProviderCount int    `json:"provider_count"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &pack); err != nil {
		return nil, err
	}
	proven := ForNetwork(config.Mainnet).Model
	out := make([]CatalogEntry, 0, len(pack.Data))
	for _, row := range pack.Data {
		id := strings.TrimSpace(row.ID)
		if id == "" {
			id = strings.TrimSpace(row.Name)
		}
		if id == "" {
			continue
		}
		e := classifyCatalog(id, row.Verifiability, proven)
		e.Name = row.Name
		e.Type = row.Type
		e.TeeAttested = row.TeeAttested
		e.ProviderCount = row.ProviderCount
		out = append(out, e)
	}
	return out, nil
}

func CatalogUsableForChat(model string, net config.Network) (ok bool, why string) {
	m := strings.TrimSpace(model)
	if m == "" || m == "host-parsed" {
		return true, "host-parsed"
	}
	sku := ForNetwork(net)
	if strings.EqualFold(m, sku.Model) && sku.ProvenE2EE && sku.Verifiability == "TeeML" {
		return true, "direct_teeml_research_sku"
	}
	return false, "not_direct_on_this_workspace"
}
