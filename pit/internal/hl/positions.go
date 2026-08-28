package hl

import (
	"encoding/json"
	"strings"
)

type Position struct {
	Coin          string `json:"coin"`
	Sz            string `json:"szi"`
	EntryPx       string `json:"entryPx,omitempty"`
	UnrealizedPnl string `json:"unrealizedPnl,omitempty"`
	Leverage      string `json:"leverage,omitempty"`
	MarginUsed    string `json:"marginUsed,omitempty"`
}

func atomString(raw json.RawMessage) string {
	if len(raw) == 0 || string(raw) == "null" {
		return ""
	}
	var s string
	if json.Unmarshal(raw, &s) == nil {
		return s
	}
	var n json.Number
	if json.Unmarshal(raw, &n) == nil {
		return n.String()
	}
	return strings.Trim(string(raw), `"`)
}

func ParsePositions(raw json.RawMessage) []Position {
	var st struct {
		AssetPositions []struct {
			Position struct {
				Coin          string          `json:"coin"`
				Szi           json.RawMessage `json:"szi"`
				EntryPx       json.RawMessage `json:"entryPx"`
				UnrealizedPnl json.RawMessage `json:"unrealizedPnl"`
				MarginUsed    json.RawMessage `json:"marginUsed"`
				Leverage      struct {
					Type  string          `json:"type"`
					Value json.RawMessage `json:"value"`
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
		lev := atomString(row.Position.Leverage.Value)
		if lev == "" {
			lev = row.Position.Leverage.Type
		}
		out = append(out, Position{
			Coin:          row.Position.Coin,
			Sz:            atomString(row.Position.Szi),
			EntryPx:       atomString(row.Position.EntryPx),
			UnrealizedPnl: atomString(row.Position.UnrealizedPnl),
			Leverage:      lev,
			MarginUsed:    atomString(row.Position.MarginUsed),
		})
	}
	return out
}

type ClearinghouseSummary struct {
	AccountValue    string `json:"accountValue,omitempty"`
	TotalMarginUsed string `json:"totalMarginUsed,omitempty"`
	TotalNtlPos     string `json:"totalNtlPos,omitempty"`
	Withdrawable    string `json:"withdrawable,omitempty"`
}

func ParseClearinghouse(raw json.RawMessage) ClearinghouseSummary {
	var st struct {
		MarginSummary struct {
			AccountValue    string `json:"accountValue"`
			TotalNtlPos     string `json:"totalNtlPos"`
			TotalMarginUsed string `json:"totalMarginUsed"`
		} `json:"marginSummary"`
		Withdrawable string `json:"withdrawable"`
	}
	if json.Unmarshal(raw, &st) != nil {
		return ClearinghouseSummary{}
	}
	return ClearinghouseSummary{
		AccountValue:    st.MarginSummary.AccountValue,
		TotalMarginUsed: st.MarginSummary.TotalMarginUsed,
		TotalNtlPos:     st.MarginSummary.TotalNtlPos,
		Withdrawable:    st.Withdrawable,
	}
}

func (c *Client) Positions(user string) ([]Position, error) {
	rows, _, err := c.Clearinghouse(user)
	return rows, err
}

func (c *Client) Clearinghouse(user string) ([]Position, ClearinghouseSummary, error) {
	raw, err := c.postInfo(map[string]any{"type": "clearinghouseState", "user": user})
	if err != nil {
		return nil, ClearinghouseSummary{}, err
	}
	return ParsePositions(raw), ParseClearinghouse(raw), nil
}
