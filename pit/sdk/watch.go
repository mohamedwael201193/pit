package sdk

func (c Client) WatchNeverPlaces() bool {
	return !c.WatchMayTrade()
}
