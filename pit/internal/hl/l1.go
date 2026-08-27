package hl

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/vmihailenco/msgpack/v5"
)

type packOrder struct {
	Type     string          `msgpack:"type"`
	Orders   []packOrderItem `msgpack:"orders"`
	Grouping string          `msgpack:"grouping"`
}

type packOrderItem struct {
	A int        `msgpack:"a"`
	B bool       `msgpack:"b"`
	P string     `msgpack:"p"`
	S string     `msgpack:"s"`
	R bool       `msgpack:"r"`
	T packOrderT `msgpack:"t"`
	C string     `msgpack:"c"`
}

type packOrderT struct {
	Limit packLimit `msgpack:"limit"`
}

type packLimit struct {
	Tif string `msgpack:"tif"`
}

func packAction(raw json.RawMessage) ([]byte, error) {
	var head struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(raw, &head); err != nil {
		return nil, err
	}
	if err := AssertActionType(raw); err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	enc := msgpack.NewEncoder(&buf)
	enc.UseCompactInts(true)
	switch head.Type {
	case "order":
		var w OrderWire
		if err := json.Unmarshal(raw, &w); err != nil {
			return nil, err
		}
		p := packOrder{Type: w.Type, Grouping: w.Grouping}
		for _, it := range w.Orders {
			p.Orders = append(p.Orders, packOrderItem{
				A: it.A, B: it.B, P: it.P, S: it.S, R: it.R,
				T: packOrderT{Limit: packLimit{Tif: it.T.Limit.Tif}},
				C: it.C,
			})
		}
		if err := enc.Encode(p); err != nil {
			return nil, err
		}
	case "cancel":
		var w CancelWire
		if err := json.Unmarshal(raw, &w); err != nil {
			return nil, err
		}
		if err := enc.Encode(w); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsigned_action")
	}
	return buf.Bytes(), nil
}

func ActionHash(raw json.RawMessage, nonce uint64) ([32]byte, error) {
	packed, err := packAction(raw)
	if err != nil {
		return [32]byte{}, err
	}
	var nonceB [8]byte
	binary.BigEndian.PutUint64(nonceB[:], nonce)
	data := append(packed, nonceB[:]...)
	data = append(data, 0x00) // no vault
	sum := crypto.Keccak256(data)
	var out [32]byte
	copy(out[:], sum)
	return out, nil
}
