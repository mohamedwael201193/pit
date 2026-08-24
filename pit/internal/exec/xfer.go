package exec

func RefuseTransfer(action string) error {
	return RefuseWithdraw(action)
}
