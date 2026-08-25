package market

func OptionalPolymarket(q Quote) error {
	if q.Source != "polymarket" {
		return nil
	}
	return Validate(q)
}
