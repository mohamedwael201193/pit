package hl

import (
	"encoding/json"
	"fmt"
)

type Trade struct {
	Coin string  `json:"coin"`
	Px   float64 `json:"px"`
	Sz   float64 `json:"sz"`
	Time int64   `json:"time"`
	Side string  `json:"side"`
}

func ParseTrades(raw json.RawMessage) ([]Trade, error) {
	var rows []map[string]any
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("no_trades")
	}
	out := make([]Trade, 0, len(rows))
	for _, r := range rows {
		t := Trade{
			Coin: fmt.Sprint(r["coin"]),
			Px:   asFloat(r["px"]),
			Sz:   asFloat(r["sz"]),
			Time: int64(asFloat(r["time"])),
			Side: fmt.Sprint(r["side"]),
		}
		if t.Px <= 0 || t.Time == 0 {
			return nil, fmt.Errorf("bad_trade")
		}
		out = append(out, t)
	}
	return out, nil
}
