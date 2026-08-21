package hl

import (
	"encoding/json"
	"fmt"
)

type Impact struct {
	Coin string
	Buy  float64
	Sell float64
}

func ParseImpactPxs(raw json.RawMessage) (Impact, error) {
	var pxs []any
	if err := json.Unmarshal(raw, &pxs); err != nil {
		return Impact{}, err
	}
	if len(pxs) < 2 {
		return Impact{}, fmt.Errorf("impact_shape")
	}
	return Impact{Buy: asFloat(pxs[0]), Sell: asFloat(pxs[1])}, nil
}
