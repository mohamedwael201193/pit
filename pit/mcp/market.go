package mcp

import "github.com/mohamedwael201193/pit/internal/market"

func MarketQuote(q market.Quote) Response {
	if err := market.Validate(q); err != nil {
		return Response{OK: false, Error: err.Error()}
	}
	return Response{OK: true, Body: q}
}
