package ledger

import "fmt"

func AfterTimeout(exchangeKnown bool) error {
	_, err := Recover(Record{Status: StatusTimeout}, "", exchangeKnown)
	if err != nil {
		return err
	}
	if exchangeKnown {
		return nil
	}
	return fmt.Errorf("query_exchange_do_not_repost")
}
