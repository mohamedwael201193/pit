package sdk

import (
	"github.com/mohamedwael201193/pit/internal/chain8004"
	"github.com/mohamedwael201193/pit/internal/identity"
)

func (c Client) FeedbackAllowed(owner, reporter, caller string) error {
	return chain8004.FeedbackAllowed(identity.Address(owner), identity.Address(reporter), identity.Address(caller))
}
