package compute

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// ProviderModel is one row from a Direct provider's /v1/models. Not the Router catalog.
type ProviderModel struct {
	ID            string
	Verifiability string
}

func FetchProviderModels(rawURL string) ([]ProviderModel, error) {
	if err := DenyRouter(rawURL); err != nil {
		return nil, err
	}
	base := strings.TrimRight(strings.TrimSpace(rawURL), "/")
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, base+"/v1/models", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "pit/1.0")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("provider_models_http_%d", resp.StatusCode)
	}
	var pack struct {
		Data []struct {
			ID            string `json:"id"`
			Verifiability string `json:"verifiability"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &pack); err != nil {
		return nil, err
	}
	out := make([]ProviderModel, 0, len(pack.Data))
	for _, row := range pack.Data {
		id := strings.TrimSpace(row.ID)
		if id == "" {
			continue
		}
		out = append(out, ProviderModel{ID: id, Verifiability: strings.TrimSpace(row.Verifiability)})
	}
	return out, nil
}
