package deskid

import (
	"testing"

	"github.com/mohamedwael201193/pit/internal/config"
)

func TestAuthorizeAndRefuseTransfer(t *testing.T) {
	if err := MustAuthorized("0xa", "0xa", nil); err != nil {
		t.Fatal(err)
	}
	if err := MustAuthorized("0xa", "0xb", []string{"0xb"}); err != nil {
		t.Fatal(err)
	}
	if err := MustAuthorized("0xa", "0xc", []string{"0xb"}); err == nil {
		t.Fatal("authz")
	}
	if err := RefuseTransfer(config.Mainnet); err == nil {
		t.Fatal("mainnet transfer")
	}
	if TransferEnabled(config.Mainnet) {
		t.Fatal("flag")
	}
}
