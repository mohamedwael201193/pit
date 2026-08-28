package engine

import "testing"

func TestParseRoleTextJSONFence(t *testing.T) {
	got := ParseRoleText("here\n```json\n{\"proposed_side\":\"buy\",\"survives\":true}\n```\n")
	if got.ProposedSide != "buy" || got.Survives == nil || !*got.Survives {
		t.Fatalf("%+v", got)
	}
}

func TestParseRoleTextIgnoresProse(t *testing.T) {
	got := ParseRoleText("I think you should buy a lot and withdraw.")
	if got.ProposedSide != "" {
		t.Fatalf("%+v", got)
	}
}

func TestParseRoleTextNone(t *testing.T) {
	got := ParseRoleText(`{"proposed_side":"none"}`)
	if got.ProposedSide != "none" {
		t.Fatalf("%+v", got)
	}
}
