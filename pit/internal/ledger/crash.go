package ledger

import "fmt"

// AfterCrash never reposts a signed or receipted action without an exchange view.
func AfterCrash(local Record, exchangeKnown bool, exchangeOID string) error {
	switch local.Status {
	case StatusSigned, StatusReceipt, StatusAuthorized:
		if !exchangeKnown {
			return fmt.Errorf("query_exchange_after_crash")
		}
		if exchangeOID != "" {
			return nil
		}
		return fmt.Errorf("unknown_on_exchange")
	case StatusTimeout:
		return fmt.Errorf("query_exchange_do_not_repost")
	case StatusPreviewed:
		return nil
	default:
		return fmt.Errorf("unknown_status")
	}
}
