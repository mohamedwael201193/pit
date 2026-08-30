package mcp

import "testing"

func TestTradeDenied(t *testing.T) {
	if !TradeDenied("order") || !TradeDenied("authorize") {
		t.Fatal("order")
	}
	if TradeDenied("opportunities") || TradeDenied("market") {
		t.Fatal("read")
	}
	if !TradeDenied("sealer") || !TradeDenied("post") || !TradeDenied("mission") {
		t.Fatal("forbidden")
	}
	if Handle(Request{Tool: "order"}).OK {
		t.Fatal("mcp order")
	}
}
