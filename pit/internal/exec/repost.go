package exec

import (
	"fmt"

	"github.com/mohamedwael201193/pit/internal/ledger"
)

// AfterTimeout decides the next step. It never blindly reposts a signed or timed-out order.
func AfterTimeout(local ledger.Record, exchangeOID string, exchangeKnown bool) error {
	ok, err := ledger.Recover(local, exchangeOID, exchangeKnown)
	if err != nil {
		return err
	}
	if ok && local.Status == ledger.StatusPreviewed {
		return nil
	}
	if !ok {
		return fmt.Errorf("do_not_repost")
	}
	return fmt.Errorf("do_not_repost")
}
