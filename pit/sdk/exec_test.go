package sdk

import "testing"

func TestSDKCannotExecute(t *testing.T) {
	c := Client{Network: "mainnet"}
	if c.CanExecute() || BrowserCanSign() {
		t.Fatal("exec")
	}
}
