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
		t.Fatal("unacked glm-5.2 twin")
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
	if err := MatchFrozenSKU(got, want); err != nil {
		t.Fatalf("live getService drifted from frozen Direct glm-5.2: %v %+v", err, got)
	}
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
