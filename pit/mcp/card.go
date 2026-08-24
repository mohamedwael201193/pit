package mcp

func EmptyCard() Response {
	return Response{OK: true, Body: map[string]any{
		"copy":     "Not enough resolved forecasts.",
		"invented": false,
		"n":        0,
	}}
}
