package exec

var PreviewBindFields = []string{
	"market", "side", "sz", "orderType", "limitPx", "slippageBps",
	"policyVersion", "sessionId", "expiryUnixMs", "nonce", "cloid", "forecastId",
}

func MissingBindField(present map[string]bool) string {
	for _, f := range PreviewBindFields {
		if !present[f] {
			return f
		}
	}
	return ""
}
