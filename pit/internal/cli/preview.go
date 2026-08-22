package cli

const PreviewCopy = "exact preview binds market, side, size, order type, limit, slippage, policy, session, expiry, nonce, clientOrderId, forecast"

func MutationInvalidates() string {
	return "any mutation invalidates authorization"
}
