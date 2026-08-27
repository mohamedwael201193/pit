package cli

import (
	"testing"

	"github.com/mohamedwael201193/pit/internal/hl"
)

func TestPostLinkedRefusesUnlinked(t *testing.T) {
	if _, err := PostLinked("testnet", hl.Envelope{}, false, "h"); err == nil {
		t.Fatal("unlinked")
	}
}
