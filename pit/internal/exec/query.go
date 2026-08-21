package exec

import "fmt"

// QueryBeforeRetry refuses a second post until the venue is queried.
func QueryBeforeRetry(exchangeKnown bool, localOID, exchangeOID string) error {
	if !exchangeKnown {
		return fmt.Errorf("query_exchange_first")
	}
	if localOID != "" && exchangeOID != "" && stringsEqual(localOID, exchangeOID) {
		return fmt.Errorf("already_posted")
	}
	return fmt.Errorf("do_not_repost")
}

func stringsEqual(a, b string) bool {
	return a == b
}
