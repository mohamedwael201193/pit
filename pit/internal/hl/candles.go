package hl

import (
	"encoding/json"
	"time"
)

func (c *Client) Candles(coin, interval string, lookback time.Duration) []map[string]any {
	if lookback <= 0 {
		lookback = 24 * time.Hour
	}
	end := time.Now().UnixMilli()
	start := time.Now().Add(-lookback).UnixMilli()
	raw, err := c.postInfo(map[string]any{
		"type": "candleSnapshot",
		"req": map[string]any{
			"coin":      coin,
			"interval":  interval,
			"startTime": start,
			"endTime":   end,
		},
	})
	if err != nil {
		return nil
	}
	var rows []map[string]any
	if json.Unmarshal(raw, &rows) != nil {
		return nil
	}
	return rows
}
