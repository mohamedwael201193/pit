package policy

import "testing"

func TestHaltDaily(t *testing.T) {
	p := Default()
	p.DailyLossUSD = 50
	if err := HaltDaily(p, -50); err == nil {
		t.Fatal("halt")
	}
	if err := HaltDaily(p, -10); err != nil {
		t.Fatal(err)
	}
}
