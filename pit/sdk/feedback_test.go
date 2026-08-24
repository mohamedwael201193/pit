package sdk

import (
	"testing"

	"github.com/mohamedwael201193/pit/internal/identity"
)

func TestFeedbackOwnerBlocked(t *testing.T) {
	o, _ := identity.NormalizeAddress("0x1111111111111111111111111111111111111111")
	r, _ := identity.NormalizeAddress("0x2222222222222222222222222222222222222222")
	c := Client{}
	if err := c.FeedbackAllowed(string(o), string(r), string(o)); err == nil {
		t.Fatal("owner")
	}
}
