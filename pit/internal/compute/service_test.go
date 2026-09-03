package compute

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mohamedwael201193/pit/internal/config"
)

func TestApplyDecodedService(t *testing.T) {
	want := MainnetChat()
	var got LiveService
	applyDecodedService([]any{
		common.HexToAddress(want.Provider),
		"chatbot",
		want.URL,
		nil, nil, nil,
		want.Model,
		"TeeML",
		"",
		common.HexToAddress(want.TeeSigner),
		true,
	}, &got)
	if err := MatchFrozenSKU(got, want); err != nil {
		t.Fatal(err)
	}
}

func TestMatchFrozenSKUNeverAutoSwaps(t *testing.T) {
	want := MainnetChat()
	ok := LiveService{
		Provider: want.Provider, URL: want.URL, Model: want.Model,
		Verifiability: "TeeML", TeeSigner: want.TeeSigner, TeeAck: true, Present: true,
	}
	if err := MatchFrozenSKU(ok, want); err != nil {
		t.Fatal(err)
	}
	unacked := ok
	unacked.TeeAck = false
	unacked.URL = "https://compute-network-28.integratenetwork.work"
	if err := MatchFrozenSKU(unacked, want); err == nil {
		t.Fatal("unacked twin")
	}
	drift := ok
	drift.URL = "https://compute-network-28.integratenetwork.work"
	if err := MatchFrozenSKU(drift, want); err == nil {
		t.Fatal("url drift")
	}
	fp8 := ok
	fp8.Model = "zai-org/GLM-5-FP8"
	if err := MatchFrozenSKU(fp8, want); err == nil {
		t.Fatal("model drift")
	}
	tls := ok
	tls.Verifiability = "TeeTLS"
	if err := MatchFrozenSKU(tls, want); err == nil {
		t.Fatal("teetls")
	}
	if want.URL == "https://compute-network-28.integratenetwork.work" {
		t.Fatal("frozen sku must not be the unacked twin")
	}
	if TestnetChat().ProvenE2EE {
		t.Fatal("galileo sealed ask stays off")
	}
}

func TestLiveGetServiceMatchesMainnetChat(t *testing.T) {
	want := MainnetChat()
	got, err := GetService(config.MainnetChain(), want.Provider)
	if err != nil {
		t.Skip(err)
	}
	if err := MatchDirectPin(got, want); err != nil {
		t.Fatalf("live getService drifted from Direct pin: %v %+v", err, got)
	}
	models, merr := FetchProviderModels(want.URL)
	if merr != nil {
		t.Log("provider /v1/models", merr)
	}
	model, err := PickDirectModel(got, models)
	if err != nil {
		t.Fatalf("pick: %v live=%+v models=%+v", err, got, models)
	}
	if !strings.EqualFold(model, "glm-5.3") || !strings.EqualFold(got.Model, "glm-5.3") {
		t.Fatalf("expected glm-5.3 TeeML pin, picked=%s getService=%s", model, got.Model)
	}
	t.Logf("direct_model=%s getService.model=%s teeml_rows=%d teeSigner=%s", model, got.Model, len(models), got.TeeSigner)
	for _, row := range models {
		if strings.EqualFold(row.ID, "glm-5.3") && strings.EqualFold(row.Verifiability, "TeeTLS") && model == "glm-5.3" {
			t.Fatal("picked TeeTLS glm-5.3")
		}
	}
}

func TestLiveListServicesHasPinnedTeeML(t *testing.T) {
	want := MainnetChat()
	var rows []LiveService
	total := int64(0)
	for off := int64(0); ; off += 50 {
		chunk, n, err := ListServices(config.MainnetChain(), off, 50)
		if err != nil {
			t.Skip(err)
		}
		total = n
		rows = append(rows, chunk...)
		if off+50 >= n || len(chunk) == 0 {
			break
		}
	}
	if total < 1 || len(rows) < 1 {
		t.Fatalf("services %d total %d", len(rows), total)
	}
	found := false
	glm53TeeML := 0
	for _, row := range rows {
		t.Logf("svc model=%s ver=%s ack=%v url=%s provider=%s", row.Model, row.Verifiability, row.TeeAck, row.URL, row.Provider)
		if addrEq(row.Provider, want.Provider) && strings.EqualFold(row.Verifiability, "TeeML") && row.TeeAck {
			found = true
		}
		if strings.EqualFold(row.Model, "glm-5.3") && strings.EqualFold(row.Verifiability, "TeeML") {
			glm53TeeML++
		}
	}
	if !found {
		t.Fatal("pinned Direct provider missing from getAllServices TeeML")
	}
	t.Logf("onchain_glm53_teeml=%d listed=%d total=%d", glm53TeeML, len(rows), total)
}

func TestLivePubkeyMatchesFrozenTeeSigner(t *testing.T) {
	want := MainnetChat()
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	signer, err := fetchPubkeySigner(ctx, want.URL+"/v1/e2ee/pubkey")
	if err != nil {
		t.Skip(err)
	}
	if !addrEq(signer, want.TeeSigner) {
		t.Fatalf("pubkey signer %s != frozen %s", signer, want.TeeSigner)
	}
}

func TestLiveRouterCatalogNeverPrivateBook(t *testing.T) {
	rows, err := ListOfficialCatalog(context.Background())
	if err != nil {
		t.Skip(err)
	}
	if len(rows) == 0 {
		t.Fatal("empty catalog")
	}
	teeml := 0
	for _, row := range rows {
		if row.PrivateBook {
			t.Fatalf("router listing marked private book: %s", row.ID)
		}
		if strings.EqualFold(row.Verifiability, "TeeML") {
			teeml++
		}
	}
	if teeml < 5 {
		t.Fatalf("router TeeML %d", teeml)
	}
}
