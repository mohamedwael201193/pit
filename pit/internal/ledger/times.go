package ledger

func AfterUnknownExchange(localStatus string, exchangeKnown bool) string {
	if !exchangeKnown {
		return "query_first"
	}
	return localStatus
}
