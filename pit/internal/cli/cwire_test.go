package cli

import (
	"testing"

	"github.com/mohamedwael201193/pit/internal/hl"
)

func TestCancelWireIsCancelOnly(t *testing.T) {
	raw, err := CancelWire(1, "0x11111111111111111111111111111111")
	if err != nil {
		t.Fatal(err)
	}
	if err := hl.AssertActionType(raw); err != nil {
		t.Fatal(err)
	}
}
