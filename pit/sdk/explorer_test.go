package sdk

import "testing"

func TestRefuseWrongExplorer(t *testing.T) {
	c := Client{Network: "mainnet"}
	if err := c.RefuseWrongExplorer("https://chainscan.0g.ai/tx/0x1"); err != nil {
		t.Fatal(err)
	}
	if err := c.RefuseWrongExplorer("https://chainscan-galileo.0g.ai/tx/0x1"); err == nil {
		t.Fatal("mix")
	}
}
