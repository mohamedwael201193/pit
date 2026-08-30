package mcp

func PreparePreview() Response {
	return Response{OK: true, Body: map[string]any{
		"prepare":   false,
		"authorize": false,
		"sign":      false,
		"trade":     false,
		"note":      "MCP cannot prepare or authorize a preview. Open PIT Desktop or CLI and type AUTHORIZE there.",
	}}
}
