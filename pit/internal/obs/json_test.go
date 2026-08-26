package obs

import (
	"testing"

	"github.com/mohamedwael201193/pit/internal/watch"
)

func TestWatchBodyNeverSigns(t *testing.T) {
	body := WatchBody(watch.EmptyPublic("mainnet"))
	if err := RefuseHealthSecrets(body); err != nil {
		t.Fatal(err)
	}
	if body["sign"] != false || body["trade"] != false {
		t.Fatalf("%+v", body)
	}
}

func TestHealthBodyNoWallet(t *testing.T) {
	b := HealthBody("req-1")
	if err := RefuseHealthSecrets(b); err != nil {
		t.Fatal(err)
	}
	if b["sign"] != false {
		t.Fatal("sign")
	}
	if _, ok := b["wallet"]; ok {
		t.Fatal("wallet")
	}
}

func TestRefuseHealthSecretsBook(t *testing.T) {
	if err := RefuseHealthSecrets(map[string]any{"private_book": true}); err == nil {
		t.Fatal("book")
	}
}
