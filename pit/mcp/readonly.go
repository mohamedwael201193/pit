package mcp

func TradeDenied(tool string) bool {
	switch tool {
	case "authorize", "order", "cancel", "transfer", "withdraw", "export_session", "key":
		return true
	default:
		return false
	}
}
