package verify

func RefuseTestnetReceiptOnMainnet(receiptNetwork, workspaceNetwork string) error {
	return SameNetwork(receiptNetwork, workspaceNetwork)
}
