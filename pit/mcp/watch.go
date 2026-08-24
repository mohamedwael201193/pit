package mcp

func WatchNeverTrades() bool {
	r := Handle(Request{Tool: "opportunities"})
	body, _ := r.Body.(map[string]any)
	trade, _ := body["trade"].(bool)
	return r.OK && !trade
}
