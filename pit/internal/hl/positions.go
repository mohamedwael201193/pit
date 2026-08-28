package hl

import "encoding/json"

type Position struct {
	Coin          string `json:"coin"`
	Sz            string `json:"szi"`
	EntryPx       string `json:"entryPx,omitempty"`
	UnrealizedPnl string `json:"unrealizedPnl,omitempty"`
	Leverage      string `json:"leverage,omitempty"`
	MarginUsed    string `json:"marginUsed,omitempty"`
}

func ParsePositions(raw json.RawMessage) []Position {
	var st struct {
		AssetPositions []struct {
			Position struct {
				Coin          string `json:"coin"`
				Szi           string `json:"szi"`
				EntryPx       string `json:"entryPx"`
				UnrealizedPnl string `json:"unrealizedPnl"`
				MarginUsed    string `json:"marginUsed"`
				Leverage      struct {
					Type  string `json:"type"`
					Value string `json:"value"`
				} `json:"leverage"`
			} `json:"position"`
		} `json:"assetPositions"`
	}
	if json.Unmarshal(raw, &st) != nil {
		return nil
	}
	out := make([]Position, 0, len(st.AssetPositions))
	for _, row := range st.AssetPositions {
		if row.Position.Coin == "" {
			continue
		}
		lev := row.Position.Leverage.Value
		if lev == "" {
			lev = row.Position.Leverage.Type
		}
		out = append(out, Position{
			Coin:          row.Position.Coin,
			Sz:            row.Position.Szi,
			EntryPx:       row.Position.EntryPx,
			UnrealizedPnl: row.Position.UnrealizedPnl,
			Leverage:      lev,
			MarginUsed:    row.Position.MarginUsed,
		})
	}
	return out
}

func (c *Client) Positions(user string) ([]Position, error) {
	raw, err := c.postInfo(map[string]any{"type": "clearinghouseState", "user": user})
	if err != nil {
		return nil, err
	}
	return ParsePositions(raw), nil
}
