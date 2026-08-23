package mcp

func VerifyHint() map[string]any {
	return map[string]any{
		"sign":     false,
		"authorize": false,
		"hint":     "recompute the storage proof with the official Go client --proof",
	}
}
