package deskid

import (
	"testing"

	"github.com/mohamedwael201193/pit/internal/config"
)

func TestBeforeSealedAsk(t *testing.T) {
	if err := BeforeSealedAsk(false); err == nil {
		t.Fatal("authz")
	}
	if err := BeforeSealedAsk(true); err != nil {
		t.Fatal(err)
	}
}

func TestRefuseTransferMainnet(t *testing.T) {
	if err := RefuseTransfer(config.Mainnet); err == nil {
		t.Fatal("mainnet")
	}
}
