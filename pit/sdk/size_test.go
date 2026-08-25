package sdk

import "testing"

func TestDropModelSize(t *testing.T) {
	c := Client{Network: "mainnet"}
	m := map[string]any{"sizeUsd": 1e9, "p": 0.9}
	if err := c.DropModelSize(m); err != nil {
		t.Fatal(err)
	}
	if _, ok := m["sizeUsd"]; ok {
		t.Fatal("size")
	}
}
