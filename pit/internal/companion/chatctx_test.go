package companion

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mohamedwael201193/pit/internal/compute"
	"github.com/mohamedwael201193/pit/internal/watch"
)

func TestChatWhatIsHappeningIdle(t *testing.T) {
	h := New(t.TempDir())
	req := local(httptest.NewRequest(http.MethodPost, "/local/chat", bytes.NewBufferString(`{"text":"What's happening?"}`)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatal(rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "Idle") {
		t.Fatal(rec.Body.String())
	}
	if strings.Contains(rec.Body.String(), `"execute":true`) || strings.Contains(rec.Body.String(), `"posted":true`) {
		t.Fatal(rec.Body.String())
	}
}

func TestChatCannotExecute(t *testing.T) {
	h := New(t.TempDir())
	req := local(httptest.NewRequest(http.MethodPost, "/local/chat", bytes.NewBufferString(`{"text":"Why can't PIT execute this?"}`)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatal(rec.Body.String())
	}
	if strings.Contains(rec.Body.String(), `"execute":true`) || strings.Contains(rec.Body.String(), `"posted":true`) {
		t.Fatal(rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "AUTHORIZE") {
		t.Fatal(rec.Body.String())
	}
}

func TestFormatBookFloorsListsPerMarketShortfall(t *testing.T) {
	got := formatBookFloors(9.38, []watch.PublicCoin{
		{Coin: "AVAX", PolicyEligible: true, MinNotional: 10.02},
		{Coin: "BTC", PolicyEligible: true, MinNotional: 10.16},
		{Coin: "SOL", Eligible: false, PolicyEligible: false, MinNotional: 10.52},
	})
	if !strings.Contains(got, "AVAX $0.64 short of $10.02") {
		t.Fatal(got)
	}
	if !strings.Contains(got, "BTC $0.78 short of $10.16") {
		t.Fatal(got)
	}
	if strings.Contains(got, "SOL") {
		t.Fatal(got)
	}
	if strings.Contains(got, "Latest:") {
		t.Fatal(got)
	}
}

func TestFormatBookFloorsPolicyClipTightNotFund(t *testing.T) {
	got := formatBookFloors(16.18, []watch.PublicCoin{
		{Coin: "ETH", PolicyEligible: true, MinNotional: 10.08, PolicyClip: 10, ExecGate: "policy_clip_tight"},
	})
	if !strings.Contains(got, "policy cap $0.08 too tight") {
		t.Fatal(got)
	}
	if strings.Contains(strings.ToLower(got), "fund") || strings.Contains(got, "short of") {
		t.Fatal(got)
	}
}

func TestLocalModelsDirectOnly(t *testing.T) {
	old := fetchOfficialCatalog
	fetchOfficialCatalog = func(ctx context.Context) ([]compute.CatalogEntry, error) {
		return []compute.CatalogEntry{{ID: "claude-opus-5", Verifiability: "unproven", Path: "catalog-listing", PrivateBook: false, Note: "listing only"}}, nil
	}
	t.Cleanup(func() { fetchOfficialCatalog = old })
	h := New(t.TempDir())
	req := local(httptest.NewRequest(http.MethodGet, "/local/models", nil))
	rec := httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatal(rec.Body.String())
	}
	body := strings.ToLower(rec.Body.String())
	if !strings.Contains(body, "glm-5.2") || !strings.Contains(body, `"path":"direct"`) {
		t.Fatal(rec.Body.String())
	}
	if strings.Contains(body, "router-api") {
		t.Fatal("router inference url leaked")
	}
	var got map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	models, _ := got["models"].([]any)
	if len(models) != 1 {
		t.Fatalf("catalog %d", len(models))
	}
	row0, _ := models[0].(map[string]any)
	if row0["private_book"] != true {
		t.Fatal("direct sku must be private book")
	}
	groups, _ := got["groups"].(map[string]any)
	if groups == nil {
		t.Fatal("groups")
	}
	un, _ := groups["unsupported"].([]any)
	if len(un) == 0 {
		t.Fatal("unsupported")
	}
	row, _ := un[0].(map[string]any)
	if row["private_book"] == true {
		t.Fatal("catalog marked private")
	}
	off, _ := groups["official_catalog"].([]any)
	if len(off) == 0 {
		t.Fatal("official catalog missing")
	}
	cat, _ := off[0].(map[string]any)
	if cat["private_book"] == true || cat["path"] == "Direct" {
		t.Fatalf("%+v", cat)
	}
}
