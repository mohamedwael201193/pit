package hl

import (
	"encoding/json"

	"github.com/mohamedwael201193/pit/internal/session"
)

type OrderWire struct {
	Type    string        `json:"type"`
	Orders  []OrderItem   `json:"orders"`
	Grouping string       `json:"grouping"`
}

type OrderItem struct {
	A      int     `json:"a"`
	B      bool    `json:"b"`
	P      string  `json:"p"`
	S      string  `json:"s"`
	R      bool    `json:"r"`
	T      OrderT  `json:"t"`
	C      string  `json:"c"`
}

type OrderT struct {
	Limit struct {
		Tif string `json:"tif"`
	} `json:"limit"`
}

type CancelWire struct {
	Type    string       `json:"type"`
	Cancels []CancelItem `json:"cancels"`
}

type CancelItem struct {
	A int    `json:"a"`
	C string `json:"c"`
}

func BuildOrder(asset int, buy bool, px, sz, cloid string) (json.RawMessage, error) {
	if err := session.CheckAction("order"); err != nil {
		return nil, err
	}
	if err := ValidCloid(cloid); err != nil {
		return nil, err
	}
	w := OrderWire{Type: "order", Grouping: "na"}
	item := OrderItem{A: asset, B: buy, P: px, S: sz, R: false, C: cloid}
	item.T.Limit.Tif = "Gtc"
	w.Orders = []OrderItem{item}
	return json.Marshal(w)
}

func BuildCancel(asset int, cloid string) (json.RawMessage, error) {
	if err := session.CheckAction("cancel"); err != nil {
		return nil, err
	}
	if err := ValidCloid(cloid); err != nil {
		return nil, err
	}
	w := CancelWire{Type: "cancel", Cancels: []CancelItem{{A: asset, C: cloid}}}
	return json.Marshal(w)
}

func AssertActionType(raw json.RawMessage) error {
	var head struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(raw, &head); err != nil {
		return err
	}
	return session.CheckAction(head.Type)
}
