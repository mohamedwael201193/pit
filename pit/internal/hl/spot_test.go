package hl

import "testing"

func TestSpotCountsAsFunded(t *testing.T) {
	if !SpotCountsAsFunded(ParseAccount(0, 12.5)) {
		t.Fatal("spot")
	}
	if SpotCountsAsFunded(ParseAccount(0, 0)) {
		t.Fatal("empty")
	}
}
