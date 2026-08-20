package engine

import (
	"fmt"
	"strings"
)

type Market string

func ParseMarket(s string) (Market, error) {
	s = strings.TrimSpace(s)
	switch {
	case strings.HasPrefix(s, "hyperliquid:perp:"):
		return Market(s), nil
	case strings.HasPrefix(s, "0g:univ3:"):
		return Market(s), nil
	case s == "none":
		return Market("none"), nil
	default:
		return "", fmt.Errorf("bad_market")
	}
}

type Forecast struct {
	Market      Market  `json:"market"`
	Side        string  `json:"side"`
	Invalidation string `json:"invalidation"`
	Uncertainty float64 `json:"uncertainty"`
	// P is host-derived, never an LLM probability.
	P float64 `json:"p"`
}

func BuildForecast(market, side, invalidation string, uncertainty, hostP float64) (Forecast, error) {
	m, err := ParseMarket(market)
	if err != nil {
		return Forecast{}, err
	}
	if invalidation == "" {
		return Forecast{}, fmt.Errorf("invalidation_required")
	}
	if hostP < 0 || hostP > 1 {
		return Forecast{}, fmt.Errorf("host_p")
	}
	return Forecast{Market: m, Side: side, Invalidation: invalidation, Uncertainty: uncertainty, P: hostP}, nil
}

func IgnoreModelP(model map[string]any, hostP float64) float64 {
	_ = model["p"]
	_ = model["probability"]
	return hostP
}
