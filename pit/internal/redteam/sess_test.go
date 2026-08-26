package redteam

import "testing"

func TestSessionMetaMustNotHoldAKey(t *testing.T) {
	if SessionMetaLooksPublic(`{"session_key":"deadbeef"}`) {
		t.Fatal("key")
	}
	if !SessionMetaLooksPublic(`{"agentAddr":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"}`) {
		t.Fatal("addr")
	}
}
