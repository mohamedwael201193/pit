package exec

func RefuseApproveAgent(action string) error {
	return RefuseWithdraw(action)
}
