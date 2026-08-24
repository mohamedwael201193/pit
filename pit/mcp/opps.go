package mcp

import "github.com/mohamedwael201193/pit/internal/market"

func OpportunityQuote(q market.Quote) Response {
	if err := market.Validate(q); err != nil {
		return Response{OK: false, Error: err.Error()}
	}
	return Response{OK: true, Body: map[string]any{"quote": q, "trade": false, "count": 0}}
}
