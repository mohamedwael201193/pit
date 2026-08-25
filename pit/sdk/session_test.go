package sdk

import "testing"

func TestSessionNameNeverLeavesHost(t *testing.T) {
	c := Client{Network: "mainnet"}
	if !c.SessionNameNeverLeavesHost() || c.Status().CanSign {
		t.Fatal("session")
	}
}
