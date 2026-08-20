package exec

import "testing"

func TestExchangeRejectsMockAndWithdraw(t *testing.T) {
	e := NewExchange("https://mock.local/exchange")
	if err := e.Guard([]byte(`{"type":"order"}`)); err == nil {
		t.Fatal("mock")
	}
	e = NewExchange("https://api.hyperliquid.xyz/exchange")
	if err := e.Guard([]byte(`{"type":"withdraw3"}`)); err == nil {
		t.Fatal("withdraw")
	}
	raw, err := jsonOrder()
	if err != nil {
		t.Fatal(err)
	}
	if err := e.Guard(raw); err != nil {
		t.Fatal(err)
	}
}

func jsonOrder() ([]byte, error) {
	return []byte(`{"type":"order","orders":[{"a":1,"b":true,"p":"1","s":"1","r":false,"t":{"limit":{"tif":"Gtc"}},"c":"0x11111111111111111111111111111111"}],"grouping":"na"}`), nil
}
