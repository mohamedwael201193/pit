package ui

import "testing"

func TestRefuseSeedPrompt(t *testing.T) {
	if err := RefuseSeedPrompt("enter your seed phrase"); err == nil {
		t.Fatal("seed")
	}
	if err := RefuseSeedPrompt("PIT never asks for a seed phrase"); err != nil {
		t.Fatal(err)
	}
}
