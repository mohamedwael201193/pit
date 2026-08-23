//go:build live

package hl_test

import (
	"testing"

	"github.com/mohamedwael201193/pit/internal/config"
	"github.com/mohamedwael201193/pit/internal/hl"
)

func TestLiveETHBook(t *testing.T) {
	c := hl.New(config.MainnetChain())
	b, err := c.PublicBook("ETH")
	if err != nil {
		t.Fatal(err)
	}
	if b.MarkPx <= 0 {
		t.Fatalf("%+v", b)
	}
}

func TestLiveTestnetETHBook(t *testing.T) {
	c := hl.New(config.TestnetChain())
	b, err := c.PublicBook("ETH")
	if err != nil {
		t.Fatal(err)
	}
	if b.MarkPx <= 0 {
		t.Fatalf("%+v", b)
	}
}
