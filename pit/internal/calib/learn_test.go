package calib

import "testing"

func TestOverconfidentAndLearned(t *testing.T) {
	if !Overconfident(0.2, 0.15) {
		t.Fatal("ece")
	}
	if Overconfident(0.05, 0.15) {
		t.Fatal("ok")
	}
	before := Health{N: 10, Sufficient: false}
	after := Health{N: 30, Sufficient: true}
	if !Learned(before, after) {
		t.Fatal("learn")
	}
	if Learned(after, after) {
		t.Fatal("same")
	}
}
