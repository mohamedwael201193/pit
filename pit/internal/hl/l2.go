package hl

import (
	"encoding/json"
	"fmt"
)

type Level struct {
	Px float64 `json:"px"`
	Sz float64 `json:"sz"`
}

type L2Book struct {
	Coin   string  `json:"coin"`
	TimeMs int64   `json:"time"`
	Bids   []Level `json:"bids"`
	Asks   []Level `json:"asks"`
}

func ParseL2(raw json.RawMessage) (L2Book, error) {
	var wrap struct {
		Coin   string          `json:"coin"`
		Time   int64           `json:"time"`
		Levels [][][]any       `json:"levels"`
	}
	if err := json.Unmarshal(raw, &wrap); err != nil {
		return L2Book{}, err
	}
	if wrap.Coin == "" || len(wrap.Levels) < 2 {
		return L2Book{}, fmt.Errorf("l2_shape")
	}
	out := L2Book{Coin: wrap.Coin, TimeMs: wrap.Time}
	out.Bids = parseSide(wrap.Levels[0])
	out.Asks = parseSide(wrap.Levels[1])
	if len(out.Bids) == 0 || len(out.Asks) == 0 {
		return L2Book{}, fmt.Errorf("l2_empty")
	}
	return out, nil
}

func parseSide(rows [][]any) []Level {
	out := make([]Level, 0, len(rows))
	for _, row := range rows {
		if len(row) < 2 {
			continue
		}
		out = append(out, Level{Px: asFloat(row[0]), Sz: asFloat(row[1])})
	}
	return out
}

func BestBidAsk(book L2Book) (bid, ask float64, err error) {
	if len(book.Bids) == 0 || len(book.Asks) == 0 {
		return 0, 0, fmt.Errorf("l2_empty")
	}
	return book.Bids[0].Px, book.Asks[0].Px, nil
}
