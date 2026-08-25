package mcp

func ForecastNeverSizes() bool {
	r := Handle(Request{Tool: "forecast"})
	body, _ := r.Body.(map[string]any)
	if !r.OK {
		return false
	}
	if _, ok := body["sizeUsd"]; ok {
		return false
	}
	if _, ok := body["sz"]; ok {
		return false
	}
	return true
}
