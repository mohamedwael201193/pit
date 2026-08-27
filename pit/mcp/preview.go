package mcp

func PreparePreview() Response {
	return Response{OK: true, Body: map[string]any{
		"prepare":   true,
		"authorize": false,
		"sign":      false,
		"trade":     false,
		"note":      "MCP may prepare a preview on this machine. Type AUTHORIZE on desktop or CLI. MCP cannot authorize.",
	}}
}
