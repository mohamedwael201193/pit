package sdk

func (c Client) SessionNameNeverLeavesHost() bool {
	return !c.CanHoldSession()
}
