package chain8004

import (
	"testing"

	"github.com/mohamedwael201193/pit/internal/identity"
)

func TestStrangerCannotReport(t *testing.T) {
	o, _ := identity.NormalizeAddress("0x1111111111111111111111111111111111111111")
	r, _ := identity.NormalizeAddress("0x2222222222222222222222222222222222222222")
	x, _ := identity.NormalizeAddress("0x3333333333333333333333333333333333333333")
	if err := FeedbackAllowed(o, r, x); err == nil {
		t.Fatal("stranger")
	}
}
