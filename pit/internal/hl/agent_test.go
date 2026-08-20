package hl

import (
	"testing"
	"time"
)

func TestApproveAgentMasterOnly(t *testing.T) {
	b, err := BuildApproveAgent("0xabc", "workspace-id-here", time.Unix(1, 0).UTC())
	if err != nil {
		t.Fatal(err)
	}
	if len(b) == 0 {
		t.Fatal("empty")
	}
	if err := SessionMustNotSign("approveAgent"); err == nil {
		t.Fatal("session")
	}
	if err := SessionMustNotSign("order"); err != nil {
		t.Fatal(err)
	}
}
