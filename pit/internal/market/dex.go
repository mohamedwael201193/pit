package market

func OptionalDexScreener(q Quote) error {
	if q.Source != "dexscreener" {
		return nil
	}
	return Validate(q)
}
