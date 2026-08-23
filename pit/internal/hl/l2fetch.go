package hl

import (
	"encoding/json"
	"fmt"
)

func (c *Client) L2(coin string) (L2Book, error) {
	if coin == "" {
		return L2Book{}, fmt.Errorf("coin_required")
	}
	raw, err := c.postInfo(map[string]any{"type": "l2Book", "coin": coin})
	if err != nil {
		return L2Book{}, err
	}
	b, err := ParseL2(raw)
	if err != nil {
		return L2Book{}, err
	}
	if b.Coin == "" {
		b.Coin = coin
	}
	if b.Coin != coin {
		return L2Book{}, fmt.Errorf("l2_coin")
	}
	return b, nil
}

func Mid(book L2Book) (float64, error) {
	bid, ask, err := BestBidAsk(book)
	if err != nil {
		return 0, err
	}
	if bid <= 0 || ask <= 0 {
		return 0, fmt.Errorf("l2_px")
	}
	return (bid + ask) / 2, nil
}

func MustJSON(v any) json.RawMessage {
	b, _ := json.Marshal(v)
	return b
}
