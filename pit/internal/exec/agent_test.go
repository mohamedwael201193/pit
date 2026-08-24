package exec

import "testing"

func TestRefuseApproveAgent(t *testing.T) {
	if err := RefuseApproveAgent("approveAgent"); err == nil {
		t.Fatal("agent")
	}
	if err := RefuseApproveAgent("order"); err != nil {
		t.Fatal(err)
	}
}
