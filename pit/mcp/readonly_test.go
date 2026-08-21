package mcp

import "testing"

func TestTradeDenied(t *testing.T) {
	if !TradeDenied("order") || !TradeDenied("authorize") {
		t.Fatal("order")
	}
	if TradeDenied("opportunities") || TradeDenied("market") {
		t.Fatal("read")
	}
	if Handle(Request{Tool: "order"}).OK {
		t.Fatal("mcp order")
	}
}
